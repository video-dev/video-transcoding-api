package bitmovintypes

type ColorTransfer string

const (
	ColorTransferUnspecified  ColorTransfer = "UNSPECIFIED"
	ColorTransferBT709        ColorTransfer = "BT709"
	ColorTransferGAMMA22      ColorTransfer = "GAMMA22"
	ColorTransferGAMMA28      ColorTransfer = "GAMMA28"
	ColorTransferSMPTE170M    ColorTransfer = "SMPTE170M"
	ColorTransferSMPTE240M    ColorTransfer = "SMPTE240M"
	ColorTransferLINEAR       ColorTransfer = "LINEAR"
	ColortransferLOG          ColorTransfer = "LOG"
	ColorTranfserLOG_SQRT     ColorTransfer = "LOG_SQRT"
	ColorTransferIEC61966_2_4 ColorTransfer = "IEC61966_2_4"
	ColorTransferBT1361_ECG   ColorTransfer = "BT1361_ECG"
	ColorTransferIEC61966_2_1 ColorTransfer = "IEC61966_2_1"
	ColorTransferBT2020_10    ColorTransfer = "BT2020_10"
	ColorTransferBT2020_12    ColorTransfer = "BT2020_12"
	ColorTransferSMPTE2084    ColorTransfer = "SMPTE2084"
	ColorTransferSMPTE428     ColorTransfer = "SMPTE428"
	ColorTransferARIB_STD_B67 ColorTransfer = "ARIB_STD_B67"
)
