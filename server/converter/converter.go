package converter

import (
	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/server/pb"
)

func ToPBProvider(p entity.Provider) pb.Provider {
	switch p {
	case entity.ProviderGoogle:
		return pb.Provider_Google
	case entity.ProviderTwitter:
		return pb.Provider_Twitter
	case entity.ProviderApple:
		return pb.Provider_Apple
	default:
		panic("never happen")
	}
}

func FromPBProvider(p pb.Provider) entity.Provider {
	switch p {
	case pb.Provider_Google:
		return entity.ProviderGoogle
	case pb.Provider_Twitter:
		return entity.ProviderTwitter
	case pb.Provider_Apple:
		return entity.ProviderApple
	default:
		panic("never happen")
	}
}

func ToPBAuthentication(a *entity.Authentication) *pb.Authentication {
	return &pb.Authentication{
		Provider: ToPBProvider(a.Provider),
		SocialId: a.SocialID,
	}
}

func ToPBAuthentications(sa []*entity.Authentication) []*pb.Authentication {
	ret := make([]*pb.Authentication, len(sa))
	for i, a := range sa {
		ret[i] = ToPBAuthentication(a)
	}
	return ret
}
