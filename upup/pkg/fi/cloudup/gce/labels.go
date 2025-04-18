/*
Copyright 2019 The Kubernetes Authors.

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

package gce

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"k8s.io/kops/pkg/apis/kops"
)

const (
	// The tag name we use to differentiate multiple logically independent clusters running in the same region
	gceLabelNameKubernetesCluster = "k8s-io-cluster-name"

	GceLabelNameInstanceGroup     = "k8s-io-instance-group"
	GceLabelNameRolePrefix        = "k8s-io-role-"
	GceLabelNameEtcdClusterPrefix = "k8s-io-etcd-"
	ControlPlane                  = "control-plane"
	Bastion                       = "bastion"
	Node                          = "node"
)

// EncodeGCELabel encodes a string into an RFC1035 compatible value, suitable for use as GCE label key or value
// We use a URI inspired escaping, but with - instead of %.
func EncodeGCELabel(s string) string {
	var b bytes.Buffer

	for i := range len(s) {
		c := s[i]
		if ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') {
			b.WriteByte(c)
		} else {
			b.WriteByte('-')
			b.WriteByte("0123456789abcdef"[c>>4])
			b.WriteByte("0123456789abcdef"[c&15])
		}
	}

	return b.String()
}

// DecodeGCELabel reverse EncodeGCELabel, taking the encoded RFC1035 compatible value back to a string
func DecodeGCELabel(s string) (string, error) {
	uriForm := strings.Replace(s, "-", "%", -1)
	v, err := url.QueryUnescape(uriForm)
	if err != nil {
		return "", fmt.Errorf("cannot decode GCE label: %q", s)
	}
	return v, nil
}

// TagForRole return the instance (network) tag used for instances with the given role.
func TagForRole(clusterName string, role kops.InstanceGroupRole) string {
	return ClusterPrefixedName(GceLabelNameRolePrefix+role.ToLowerString(), clusterName, 63)
}
