package customer

//func PutImageFile(ctx context.Context, name string, rdr io.Reader) error {
//
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
//
//	writer := client.Bucket(gcsBucket).Object(name).NewWriter(ctx)
//
//	io.Copy(writer, rdr)
//	// check for errors on io.Copy in production code!
//	return writer.Close()
//}

//func PutImageDataToSql(imageFromUi ImageFromUi) error{
//
//		_, err := config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,name) VALUES (?,?,?,?)",
//			imageFromUi.SiteId, imageFromUi.PageId, imageFromUi.AltText, imageFromUi.FileName)
//
//		if err != nil {
//			//return pages, errors.New("500. Internal Server Error." + err.Error())
//			log.Fatalf("Could not INSERT into image: %v", err)
//		}
//
//
//	return nil
//}