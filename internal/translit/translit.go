package translit

import (
	tr "github.com/gen1us2k/go-translit"
	"github.com/supme/service/proto"
	"io"
)

type TrServer struct {
}

func NewTr() *TrServer {
	return &TrServer{}
}

func (srv *TrServer) EnRu(inStream proto.Transliteration_EnRuServer) error {
	for {
		inWord, err := inStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := &proto.Word{
			Word: tr.Translit(inWord.Word),
		}
		inStream.Send(out)
	}
	return nil
}
