# deSEC Webhook for cert-manager

A [cert-manager][1] webhook to solve an ACME DNS01 challenge using the [deSEC][2] API

## Prerequisites

A Kubernetes cluster with cert-manager deployed. If you haven't already installed cert-manger, follow the guide [here][3].

## Deployment

### Using regular manifests

An example webhook deployment, with it's associated roles, role bindings, service and required tls certificates, is provided in the file `examples/deploy/desec-webhook.yaml`. The example manifest will deploy to the default namespace. To set a different namespace, replace all instances of `default` with the namespace of your choice. Don't miss the annotation ` cert-manager.io/inject-ca-from` in the `APIService` or the `dnsNames` for the webhook `Certificate`. The following will replace all the necessary namespace references in the example manifest:

```bash
$ sed -i 's/default/yournamespace/g' examples/deploy/desec-webhook.yaml
```

Once you're satisfied with the changes, deploy the webhook with:

```bash
$ kubectl apply -f examples/deploy/desec-webhook.yaml
```

### Using Helm

```bash
$ helm upgrade --install desec-webhook ./deploy/desec-webhook \
  --namespace cert-manager \
  --set groupName=acme.example.com \
  --set certManager.namespace=cert-manager \
  --set certManager.serviceAccountName=cert-manager
```

If you want to scope secret access, set `secretsNames`:

```bash
$ helm upgrade --install desec-webhook ./deploy/desec-webhook \
  --namespace cert-manager \
  --set groupName=acme.example.com \
  --set certManager.namespace=cert-manager \
  --set certManager.serviceAccountName=cert-manager \
  --set-json 'secretsNames=["desec-token"]'
```

## Usage

### Deploy an API Token Secret

The deSEC API token needs to placed into a kubernetes secret. You can use the file `examples/desec-token.yaml` as a starting point. Place your Base64-encoded API token into the manifest. This can be obtained with:

```bash
$ echo -n "your-api-token" | base64
```

### Deploy an Issuer

An example `ClusterIssuer` is provided in the `examples/letsencrypt-staging-issuer.yaml` file. It uses the letsencrypt staging server, but it can be adapted for the letsencrypt production server or other acme server. `Issuer` can be used instead of `ClusterIssuer` if you only want to issue certificates in a single namespace.

### Deploy a Certificate

An example certificate manifest is provide in `examples/test-certificate.yaml`. Edit as required for production certificates.

### Check the Certificate

Use the following check the deployed certificate, replacing _test-example-com-tls_ with the **secretName** you used in the previous step.

```bash
$ kubectl get secret test-example-com-tls -o json | jq -r '.data."tls.crt"' | base64 -d | openssl x509 -text -noout
```

## Building

```bash
$ make build
```

## Releasing Multi-Arch Images

The repository includes a GitHub Actions workflow at `.github/workflows/release-image.yaml`.

- Push a tag like `v1.2.3` to publish a multi-arch image (`linux/amd64`, `linux/arm64`) to `ghcr.io/<owner>/cert-manager-webhook-desec`.
- You can also run the workflow manually from the Actions tab (`workflow_dispatch`).

### Running the test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

See [here][4] for testing details.

[1]: https://github.com/cert-manager/cert-manager
[2]: https://desec.io/
[3]: https://cert-manager.io/docs/installation/kubernetes/
[4]: ./testdata/desec/README.md
