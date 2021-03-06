// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ansible

import (
	"path/filepath"

	"github.com/operator-framework/operator-sdk/internal/scaffold/input"
)

const MoleculeDefaultPrepareFile = "prepare.yml"

type MoleculeDefaultPrepare struct {
	StaticInput
}

// GetInput - gets the input
func (m *MoleculeDefaultPrepare) GetInput() (input.Input, error) {
	if m.Path == "" {
		m.Path = filepath.Join(MoleculeDefaultDir, MoleculeDefaultPrepareFile)
	}
	m.TemplateBody = moleculeDefaultPrepareAnsibleTmpl

	return m.Input, nil
}

const moleculeDefaultPrepareAnsibleTmpl = `---
- name: Prepare
  hosts: k8s
  gather_facts: no
  tasks:
    - name: Fetch the kubeconfig
      fetch:
        dest: '{{ kubeconfig_file }}'
        flat: yes
        src: /root/.kube/config

    - name: Change the kubeconfig port to the proper value
      replace:
        regexp: '8443'
        replace: "{{ lookup('env', 'KIND_PORT') }}"
        path: '{{ kubeconfig_file }}'
      delegate_to: localhost

    - name: Wait for the Kubernetes API to become available (this could take a minute)
      uri:
        url: "http://localhost:10080/kubernetes-ready"
        status_code: 200
        validate_certs: no
      register: result
      until: (result.status|default(-1)) == 200
      retries: 60
      delay: 5
`
