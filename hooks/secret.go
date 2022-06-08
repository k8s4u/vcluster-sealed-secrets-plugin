package hooks

import (
	"context"
	"crypto/rsa"
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

// var _ hook.MutateCreateVirtual = &secretHook{}
var _ hook.MutateCreatePhysical = &secretHook{}

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

// For some reason hooks does not get triggered without this
func (s *secretHook) MutateGetVirtual(ctx context.Context, obj client.Object) (client.Object, error) {
	fmt.Println("MutateGetVirtual called")
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object %v is not a secret", obj)
	}

	// fmt.Printf("MutateGetVirtual: secret data: %s\n", secret.Data)

	/*
		pubKey, err := readPubKey()
		if err != nil {
			return nil, err
		}

		ssecret2, err := ssv1alpha1.NewSealedSecret(scheme.Codecs, pubKey, secret)
		if err != nil {
			return nil, err
		}

		fmt.Printf("MutateGetVirtual: secret encrypted data: %s\n", ssecret2.Spec.EncryptedData)
	*/
	return secret, nil
}

// Protect secret by encrypting its data before creating it unless it is already encrypted
/*
func (p *secretHook) MutateCreateVirtual(ctx context.Context, obj client.Object) (client.Object, error) {
	fmt.Println("MutateCreateVirtual called")

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

	// If we are able to unseal secret it means that is already encrypted
	_, err = ssecret.Unseal(scheme.Codecs, priKey)
	if err == nil {
		fmt.Println("MutateCreateVirtual: secret looks to be already encrypted")
		return secret, nil
	}

	// On other why let's encrypt it
	pubKey, err := readPubKey()
	if err != nil {
		return nil, err
	}

	ssecret2, err := ssv1alpha1.NewSealedSecret(scheme.Codecs, pubKey, secret)
	if err != nil {
		return nil, err
	}

	newSecretData := map[string][]byte{}
	for key, value := range ssecret2.Spec.EncryptedData {
		valueBytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, err
		}
		newSecretData[key] = valueBytes
	}
	secret.Data = newSecretData

	fmt.Println("MutateCreateVirtual: Replacing secret with encrypted version")
	return secret, nil
}
*/

// Decrypt secret data before creating physical secret
func (p *secretHook) MutateCreatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	fmt.Println("MutateCreatePhysical called")
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
	if err != nil {
		fmt.Printf("MutateCreatePhysical: Secret looks to be already non-encrypted, error: %s\n", err)
		return secret, nil
	}

	secret.Data = unsealedSecret.Data

	fmt.Println("MutateCreatePhysical: Replacing secret with non-encrypted version")
	return secret, nil
}
