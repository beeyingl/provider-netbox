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

// Code generated by terrajet. DO NOT EDIT.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type AddressObservation struct {
	ID *string `json:"id,omitempty" tf:"id,omitempty"`
}

type AddressParameters struct {

	// +kubebuilder:validation:Optional
	DNSName *string `json:"dnsName,omitempty" tf:"dns_name,omitempty"`

	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty" tf:"description,omitempty"`

	// +kubebuilder:validation:Required
	IPAddress *string `json:"ipAddress" tf:"ip_address,omitempty"`

	// +kubebuilder:validation:Optional
	InterfaceID *float64 `json:"interfaceId,omitempty" tf:"interface_id,omitempty"`

	// +kubebuilder:validation:Required
	Status *string `json:"status" tf:"status,omitempty"`

	// +kubebuilder:validation:Optional
	Tags []*string `json:"tags,omitempty" tf:"tags,omitempty"`

	// +kubebuilder:validation:Optional
	TenantID *float64 `json:"tenantId,omitempty" tf:"tenant_id,omitempty"`

	// +kubebuilder:validation:Optional
	VrfID *float64 `json:"vrfId,omitempty" tf:"vrf_id,omitempty"`
}

// AddressSpec defines the desired state of Address
type AddressSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     AddressParameters `json:"forProvider"`
}

// AddressStatus defines the observed state of Address.
type AddressStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        AddressObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Address is the Schema for the Addresss API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,netboxjet}
type Address struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AddressSpec   `json:"spec"`
	Status            AddressStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AddressList contains a list of Addresss
type AddressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Address `json:"items"`
}

// Repository type metadata.
var (
	Address_Kind             = "Address"
	Address_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: Address_Kind}.String()
	Address_KindAPIVersion   = Address_Kind + "." + CRDGroupVersion.String()
	Address_GroupVersionKind = CRDGroupVersion.WithKind(Address_Kind)
)

func init() {
	SchemeBuilder.Register(&Address{}, &AddressList{})
}