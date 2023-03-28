package sbm

import (
	"fmt"
	"io"
)

// Write writes an SBM object into the stream.
func (sbm *Sbm) Write(writer io.Writer) (err error) {

	// Write the top headers.
	err = sbm.writeTopHeaders(writer)
	if err != nil {
		return err
	}

	// Write the binary array of bits with the 'NewFromBitsArray Line' at the end.
	err = sbm.writeArrayData(writer)
	if err != nil {
		return err
	}

	// Write the bottom headers.
	err = sbm.writeBottomHeaders(writer)
	if err != nil {
		return err
	}

	return nil
}

func (sbm *Sbm) writeTopHeaders(writer io.Writer) (err error) {

	// 1. Title.
	_, err = writer.Write([]byte(Header_FormatName))
	if err != nil {
		return err
	}

	// 2. Version.
	_, err = writer.Write([]byte(fmt.Sprintf(HeaderFormat_Version, sbm.format.version)))
	if err != nil {
		return err
	}

	// 3. Width.
	_, err = writer.Write([]byte(
		fmt.Sprintf(
			HeaderFormat_Width,
			sbm.pixelArray.metaData.width,
			sbm.pixelArray.metaData.header.width.topLeft,
			sbm.pixelArray.metaData.header.width.topRight,
		),
	))
	if err != nil {
		return err
	}

	// 4. Height
	_, err = writer.Write([]byte(
		fmt.Sprintf(
			HeaderFormat_Height,
			sbm.pixelArray.metaData.height,
			sbm.pixelArray.metaData.header.height.topLeft,
			sbm.pixelArray.metaData.header.height.topRight,
		),
	))
	if err != nil {
		return err
	}

	// 5. Area.
	_, err = writer.Write([]byte(
		fmt.Sprintf(
			HeaderFormat_Area,
			sbm.pixelArray.metaData.area,
			sbm.pixelArray.metaData.header.area.topLeft,
			sbm.pixelArray.metaData.header.area.topRight,
		),
	))
	if err != nil {
		return err
	}

	return nil
}

func (sbm *Sbm) writeArrayData(writer io.Writer) (err error) {

	// Write the binary Array of Bits with the 'NewFromBitsArray Line' at the End.
	_, err = writer.Write(sbm.pixelArray.data.bytes)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(NL))
	if err != nil {
		return err
	}

	return nil
}

func (sbm *Sbm) writeBottomHeaders(writer io.Writer) (err error) {

	// 1. Width.
	_, err = writer.Write([]byte(
		fmt.Sprintf(
			HeaderFormat_Width,
			sbm.pixelArray.metaData.width,
			sbm.pixelArray.metaData.header.width.bottomLeft,
			sbm.pixelArray.metaData.header.width.bottomRight,
		),
	))
	if err != nil {
		return err
	}

	// 2. Height
	_, err = writer.Write([]byte(
		fmt.Sprintf(
			HeaderFormat_Height,
			sbm.pixelArray.metaData.height,
			sbm.pixelArray.metaData.header.height.bottomLeft,
			sbm.pixelArray.metaData.header.height.bottomRight,
		),
	))
	if err != nil {
		return err
	}

	// 3. Area.
	_, err = writer.Write([]byte(
		fmt.Sprintf(
			HeaderFormat_Area,
			sbm.pixelArray.metaData.area,
			sbm.pixelArray.metaData.header.area.bottomLeft,
			sbm.pixelArray.metaData.header.area.bottomRight,
		),
	))
	if err != nil {
		return err
	}

	return nil
}
