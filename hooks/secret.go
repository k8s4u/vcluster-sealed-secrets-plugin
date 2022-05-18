package hooks

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/loft-sh/vcluster-sdk/hook"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes/scheme"

	ssv1alpha1 "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealed-secrets/v1alpha1"
	"k8s.io/client-go/util/keyutil"
)

const (
	tlsCertFile = `/tls/tls.crt`
	tlsKeyFile  = `/tls/tls.key`
)

var (
	// ErrPrivateKeyNotRSA is returned when the private key is not a valid RSA key.
	ErrPrivateKeyNotRSA = errors.New("Private key is not an RSA key")
)

func NewSecretHook() hook.ClientHook {
	return &secretHook{}
}

// Purpose of this hook is ...
type secretHook struct{}

func (s *secretHook) Name() string {
	return "secret-hook"
}

func (s *secretHook) Resource() client.Object {
	return &corev1.Secret{}
}

var _ hook.MutateGetVirtual = &secretHook{}
var _ hook.MutateCreateVirtual = &secretHook{}

func readPrivKey() (*rsa.PrivateKey, error) {
	key, err := keyutil.PrivateKeyFromFile(tlsKeyFile)
	if err != nil {
		return nil, err
	}

	switch rsaKey := key.(type) {
	case *rsa.PrivateKey:
		return rsaKey, nil
	default:
		return nil, ErrPrivateKeyNotRSA
	}
}

func readPubKey() (*rsa.PublicKey, error) {
	pubKeys, err := keyutil.PublicKeysFromFile(tlsCertFile)
	if err != nil {
		return nil, err
	}

	switch rsaKey := pubKeys[0].(type) {
	case *rsa.PublicKey:
		return rsaKey, nil
	default:
		return nil, ErrPrivateKeyNotRSA
	}
}

// Protect secret by encrypting its data before send its content to client
func (s *secretHook) MutateGetVirtual(ctx context.Context, obj client.Object) (client.Object, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object %v is not a secret", obj)
	}

	pubKey, err := readPubKey()
	if err != nil {
		return nil, err
	}

	ssecret, err := ssv1alpha1.NewSealedSecret(scheme.Codecs, pubKey, secret)
	if err != nil {
		return nil, err
	}

	newSecretData := map[string][]byte{}
	for key, value := range ssecret.Spec.EncryptedData {
		valueBytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, err
		}
		newSecretData[key] = valueBytes
	}
	secret.Data = newSecretData

	return secret, nil
}

// Decrypt secret data before creating it
func (p *secretHook) MutateCreateVirtual(ctx context.Context, obj client.Object) (client.Object, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object %v is not a secret", obj)
	}

	privKey, err := readPrivKey()
	if err != nil {
		return nil, err
	}
	priKey := map[string]*rsa.PrivateKey{"key": privKey}

	ssecret := ssv1alpha1.SealedSecret{}

	newEncryptedData := map[string]string{}
	for key, value := range secret.Data {
		newEncryptedData[key] = string(value)
	}
	ssecret.Spec.EncryptedData = newEncryptedData

	unsealedSecret, err := ssecret.Unseal(scheme.Codecs, priKey)
	if err == nil {
		secret.Data = unsealedSecret.Data
	}

	return secret, nil
}
