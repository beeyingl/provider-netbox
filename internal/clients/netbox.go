/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/terrajet/pkg/terraform"

	"github.com/crossplane-contrib/provider-jet-netbox/apis/v1alpha1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal netbox credentials as JSON"

	keyServerUrl = "server_url"
	keyApiToken  = "api_token"

	envToken            = "NETBOX_API_TOKEN"
	envServerUrl        = "NETBOX_SERVER_URL"
	envAllowInsecure    = "NETBOX_ALLOW_INSECURE_HTTPS"
	envSkipVersionCheck = "NETBOX_SKIP_VERSION_CHECK"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1alpha1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1alpha1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		netboxCreds := map[string]string{}
		if err := json.Unmarshal(data, &netboxCreds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// set environment variables for sensitive provider configuration
		// Deprecated: In shared gRPC mode we do not support injecting
		// credentials via the environment variables. You should specify
		// credentials via the Terraform main.tf.json instead.
		/*ps.Env = []string{
			fmt.Sprintf("%s=%s", "HASHICUPS_USERNAME", netboxCreds["username"]),
			fmt.Sprintf("%s=%s", "HASHICUPS_PASSWORD", netboxCreds["password"]),
		}*/
		// set credentials in Terraform provider configuration
		/*ps.Configuration = map[string]interface{}{
			"username": netboxCreds["username"],
			"password": netboxCreds["password"],
		}*/

		ps.Configuration = map[string]interface{}{}
		if v, ok := netboxCreds[keyServerUrl]; ok {
			ps.Configuration[keyServerUrl] = v
		}
		if v, ok := netboxCreds[keyApiToken]; ok {
			ps.Configuration[keyApiToken] = v
		}

		// set environment variables for sensitive provider configuration
		ps.Env = []string{
			fmt.Sprintf("%s=%s", envToken, netboxCreds[keyApiToken]),
			fmt.Sprintf("%s=%s", envServerUrl, netboxCreds[keyServerUrl]),
			fmt.Sprintf("%s=%t", envAllowInsecure, true),
			fmt.Sprintf("%s=%t", envSkipVersionCheck, true),
		}

		return ps, nil
	}
}
