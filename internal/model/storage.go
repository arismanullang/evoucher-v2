package model

import (
	"io"
	"net/http"

	"golang.org/x/net/context"
	st "google.golang.org/cloud/storage"
)

type saveData struct {
	c       context.Context
	r       *http.Request       //http response
	w       http.ResponseWriter //http writer
	ctx     context.Context
	cleanUp []string // cleanUp is a list of filenames that need cleaning up at the end of the saving.
	failed  bool     // failed indicates that one or more of the saving steps failed.
}

func PutFile(ctx context.Context, name string, rdr io.Reader) error {

	client, err := st.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	writer := client.Bucket(GCS_BUCKET).Object(name).NewWriter(ctx)

	io.Copy(writer, rdr)
	// check for errors on io.Copy in production code!
	return writer.Close()
}

func GetFileLink(ctx context.Context, name string) (string, error) {
	client, err := st.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	attrs, err := client.Bucket(GCS_BUCKET).Object(name).Attrs(ctx)
	if err != nil {
		return "", err
	}
	return attrs.MediaLink, nil
}
