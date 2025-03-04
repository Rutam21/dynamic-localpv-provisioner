/*
Copyright 2019 The OpenEBS Authors.

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

package app

import (
	"context"

	"github.com/pkg/errors"

	blockdeviceclaim "github.com/openebs/maya/pkg/blockdeviceclaim/v1alpha1"

	mconfig "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	klog "k8s.io/klog/v2"
)

// This function performs the preupgrade related tasks for 1.0 to 1.1
func performPreupgradeTasks(ctx context.Context, kubeClient *clientset.Clientset) error {
	return addLocalPVFinalizerOnAssociatedBDCs(ctx, kubeClient)
}

// Add localpv finalizer on the BDCs that are used by PVs provisioned from localpv provisioner
func addLocalPVFinalizerOnAssociatedBDCs(ctx context.Context, kubeClient *clientset.Clientset) error {
	// Get the list of PVs that are provisioned by device based local pv provisioner
	pvList, err := kubeClient.CoreV1().PersistentVolumes().List(
		ctx,
		metav1.ListOptions{
			LabelSelector: string(mconfig.CASTypeKey) + "=local-device",
		})
	if err != nil {
		return errors.Wrap(err, "failed to list localpv based pv(s)")
	}

	for _, pvObj := range pvList.Items {
		bdcName := "bdc-" + pvObj.Name

		bdcObj, err := blockdeviceclaim.NewKubeClient().WithNamespace(getOpenEBSNamespace()).
			Get(ctx, bdcName, metav1.GetOptions{})
		if err != nil {
			// BDCs may not exist if the PV reclaimPolicy is set
			// to 'Retain' and the BDCs have been manually removed
			// Ref: github.com/openebs/openebs/issues/3363
			// TODO: Clean this part of the code up a bit.
			klog.Warningf("failed to get bdc %v", bdcName)
			continue
		}

		// Add finalizer only if deletionTimestamp is not set
		if !bdcObj.DeletionTimestamp.IsZero() {
			continue
		}

		// Add finalizer on associated BDC
		_, err = blockdeviceclaim.BuilderForAPIObject(bdcObj).BDC.AddFinalizer(LocalPVFinalizer)
		if err != nil {
			return errors.Wrapf(err, "failed to add localpv finalizer on BDC %v",
				bdcObj.Name)
		}
	}
	return nil
}
