package converter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/server/pb"
)

func TestConverter(t *testing.T) {
	original := []*entity.Authentication{{Provider: entity.ProviderGoogle, SocialID: uuid.NewString()}, {Provider: entity.ProviderApple, SocialID: uuid.NewString()}}
	converted := ToPBAuthentications(original)

	if converted[0].Provider != pb.Provider_Google || converted[0].SocialId != original[0].SocialID {
		t.Fatalf("failed convert: %+v", converted[0])
	}

	if converted[1].Provider != pb.Provider_Apple || converted[1].SocialId != original[1].SocialID {
		t.Fatalf("failed convert: %+v", converted[1])
	}
}
