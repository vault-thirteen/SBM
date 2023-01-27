// write.go.

package sbm

import (
	"fmt"
	"io"
)

// Write writes an SBM Object into the Stream.
func (sbm *Sbm) Write(
	writer io.Writer,
) (err error) {

	// Write the top Headers.
	err = sbm.writeTopHeaders(writer)
	if err != nil {
		return
	}

	// Write the binary Array of Bits with the 'NewFromBitsArray Line' at the End.
	err = sbm.writeArrayData(writer)
	if err != nil {
		return
	}

	// Write the bottom Headers.
	err = sbm.writeBottomHeaders(writer)
	if err != nil {
		return
	}

	return
}

func (sbm *Sbm) writeTopHeaders(
	writer io.Writer,
) (err error) {

	var msg []byte

	// 1. Title.
	msg = []byte(Header_FormatName)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	// 2. Version.
	msg = []byte(fmt.Sprintf(HeaderFormat_Version, sbm.format.version))
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	// 3. Width.
	msg = []byte(
		fmt.Sprintf(
			HeaderFormat_Width,
			sbm.pixelArray.metaData.width,
			sbm.pixelArray.metaData.header.width.topLeft,
			sbm.pixelArray.metaData.header.width.topRight,
		),
	)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	// 4. Height
	msg = []byte(
		fmt.Sprintf(
			HeaderFormat_Height,
			sbm.pixelArray.metaData.height,
			sbm.pixelArray.metaData.header.height.topLeft,
			sbm.pixelArray.metaData.header.height.topRight,
		),
	)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	// 5. Area.
	msg = []byte(
		fmt.Sprintf(
			HeaderFormat_Area,
			sbm.pixelArray.metaData.area,
			sbm.pixelArray.metaData.header.area.topLeft,
			sbm.pixelArray.metaData.header.area.topRight,
		),
	)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	return
}

func (sbm *Sbm) writeArrayData(
	writer io.Writer,
) (err error) {

	var msg []byte

	// Write the binary Array of Bits with the 'NewFromBitsArray Line' at the End.
	_, err = writer.Write(sbm.pixelArray.data.bytes)
	if err != nil {
		return
	}
	msg = []byte(NL)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	return
}

func (sbm *Sbm) writeBottomHeaders(
	writer io.Writer,
) (err error) {

	var msg []byte

	// 1. Width.
	msg = []byte(
		fmt.Sprintf(
			HeaderFormat_Width,
			sbm.pixelArray.metaData.width,
			sbm.pixelArray.metaData.header.width.bottomLeft,
			sbm.pixelArray.metaData.header.width.bottomRight,
		),
	)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	// 2. Height
	msg = []byte(
		fmt.Sprintf(
			HeaderFormat_Height,
			sbm.pixelArray.metaData.height,
			sbm.pixelArray.metaData.header.height.bottomLeft,
			sbm.pixelArray.metaData.header.height.bottomRight,
		),
	)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	// 3. Area.
	msg = []byte(
		fmt.Sprintf(
			HeaderFormat_Area,
			sbm.pixelArray.metaData.area,
			sbm.pixelArray.metaData.header.area.bottomLeft,
			sbm.pixelArray.metaData.header.area.bottomRight,
		),
	)
	_, err = writer.Write(msg)
	if err != nil {
		return
	}

	return
}
