package mesh

import (
	"github.com/gitlayzer/kt-connect/pkg/kt/command/general"
	opt "github.com/gitlayzer/kt-connect/pkg/kt/command/options"
	"github.com/gitlayzer/kt-connect/pkg/kt/transmission"
	"github.com/gitlayzer/kt-connect/pkg/kt/util"
	"github.com/rs/zerolog/log"
	coreV1 "k8s.io/api/core/v1"
)

func ManualMesh(svc *coreV1.Service) error {
	meshKey, meshVersion := getVersion(opt.Get().Mesh.VersionMark)
	shadowPodName := svc.Name + util.MeshPodInfix + meshVersion
	labels := getMeshLabels(meshKey, meshVersion, svc)
	annotations := make(map[string]string)
	mirror := transmission.MirrorConfig{
		Target:      opt.Get().Mesh.MirrorTarget,
		SampleRate:  opt.Get().Mesh.MirrorSampleRate,
		RedactRules: opt.Get().Mesh.MirrorRedactRules,
		LogPath:     opt.Get().Mesh.MirrorLogPath,
	}
	if err := general.CreateShadowAndInbound(shadowPodName, opt.Get().Mesh.Expose, labels,
		annotations, general.GetTargetPorts(svc), mirror); err != nil {
		return err
	}
	log.Info().Msg("---------------------------------------------------------")
	log.Info().Msgf(" Now you can update Istio rule by label '%s=%s' ", meshKey, meshVersion)
	log.Info().Msg("---------------------------------------------------------")
	return nil
}

func getMeshLabels(meshKey, meshVersion string, svc *coreV1.Service) map[string]string {
	labels := map[string]string{}
	if svc != nil {
		for k, v := range svc.Spec.Selector {
			labels[k] = v
		}
	}
	labels[util.KtRole] = util.RoleMeshShadow
	labels[meshKey] = meshVersion
	return labels
}
