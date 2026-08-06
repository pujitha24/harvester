package pvchelper

import (
	"testing"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

func TestResolvePVCResources(t *testing.T) {
	historicalSpec := corev1.PersistentVolumeClaimSpec{
		Resources: corev1.VolumeResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("10Gi"),
			},
		},
	}

	t.Run("volume expanded after snapshot uses the snapshot's current RestoreSize", func(t *testing.T) {
		restoreSize := resource.MustParse("20Gi")
		vs := &snapshotv1.VolumeSnapshot{
			Status: &snapshotv1.VolumeSnapshotStatus{
				ReadyToUse:  ptr.To(true),
				RestoreSize: &restoreSize,
			},
		}

		got := ResolvePVCResources(historicalSpec, vs)

		gotStorage := got.Requests[corev1.ResourceStorage]
		assert.Equal(t, "20Gi", gotStorage.String())
	})

	t.Run("falls back to the historical spec when RestoreSize is not yet reported", func(t *testing.T) {
		vs := &snapshotv1.VolumeSnapshot{
			Status: &snapshotv1.VolumeSnapshotStatus{
				ReadyToUse: ptr.To(true),
			},
		}

		got := ResolvePVCResources(historicalSpec, vs)

		assert.Equal(t, historicalSpec.Resources, got)
	})

	t.Run("falls back to the historical spec when the snapshot has no status", func(t *testing.T) {
		vs := &snapshotv1.VolumeSnapshot{}

		got := ResolvePVCResources(historicalSpec, vs)

		assert.Equal(t, historicalSpec.Resources, got)
	})

	t.Run("falls back to the historical spec when the snapshot is nil", func(t *testing.T) {
		got := ResolvePVCResources(historicalSpec, nil)

		assert.Equal(t, historicalSpec.Resources, got)
	})
}
