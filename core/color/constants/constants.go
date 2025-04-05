package constants

import (
	"math"
)

// Color model conversion constants
const (
	// RGB to YUV conversion coefficients
	YUV_R = 0.299
	YUV_G = 0.587
	YUV_B = 0.114
	YUV_U = 0.615
	YUV_V = 0.51499
	YUV_W = 0.10001

	// YUV to RGB conversion coefficients
	YUV_RGB_R1 = 1.13983
	YUV_RGB_G1 = -0.39465
	YUV_RGB_G2 = -0.58060
	YUV_RGB_B1 = 2.03211

	// YUV pure color coefficients
	YUV_PureRedU   = -0.14713
	YUV_PureGreenU = -0.331
	YUV_PureGreenV = -0.419
	YUV_PureBlueU  = 0.500
	YUV_PureBlueV  = -0.081

	// YUV general case coefficients
	YUV_GeneralU1 = -0.14713
	YUV_GeneralU2 = -0.28886
	YUV_GeneralU3 = 0.436

	// YUV comparison thresholds
	YUV_Threshold = 1e-10

	// YIQ conversion coefficients
	YIQ_RGB_Y1 = 0.299
	YIQ_RGB_Y2 = 0.587
	YIQ_RGB_Y3 = 0.114
	YIQ_RGB_I1 = 0.596
	YIQ_RGB_I2 = -0.275
	YIQ_RGB_I3 = -0.321
	YIQ_RGB_Q1 = 0.212
	YIQ_RGB_Q2 = -0.523
	YIQ_RGB_Q3 = 0.311

	// YIQ to RGB conversion coefficients
	YIQ_RGB_R1 = 0.9563
	YIQ_RGB_R2 = 0.6210
	YIQ_RGB_G1 = -0.2721
	YIQ_RGB_G2 = -0.6474
	YIQ_RGB_B1 = -1.1070
	YIQ_RGB_B2 = 1.7046

	// YIQ channel ranges
	YIQ_I_Min = -0.5957
	YIQ_I_Max = 0.5957
	YIQ_Q_Min = -0.5226
	YIQ_Q_Max = 0.5226

	// YCbCr conversion coefficients (ITU-R BT.601)
	YCBCR_KR      = 0.299
	YCBCR_KG      = 0.587
	YCBCR_KB      = 0.114
	YCBCR_Divisor = 0.5

	// YCbCr channel ranges
	YCBCR_Cb_Min = -0.5
	YCBCR_Cb_Max = 0.5
	YCBCR_Cr_Min = -0.5
	YCBCR_Cr_Max = 0.5

	// Color conversion thresholds
	WhiteThreshold     = 1e-10
	NearWhiteThreshold = 0.01
	ChromaThreshold    = 1e-10
	AlphaEpsilon       = 1e-10

	// Wavelength ranges for visible spectrum
	WavelengthRedMin    = 620.0
	WavelengthYellowMin = 590.0
	WavelengthGreenMin  = 520.0
	WavelengthCyanMin   = 495.0
	WavelengthBlueMin   = 450.0
	WavelengthVioletMin = 380.0

	// Hue ranges for color models
	HueRed     = 0.0
	HueYellow  = 60.0
	HueGreen   = 120.0
	HueCyan    = 180.0
	HueBlue    = 240.0
	HueMagenta = 300.0
	HueRange   = 60.0

	// Color model conversion constants
	DegreesPerCircle = 360.0
	RadiansPerDegree = 0.017453292519943295 // π/180
	DegreesPerRadian = 57.29577951308232    // 180/π
	HueNormalization = 1.0 / 360.0
	HueOffset1       = 1.0 / 3.0
	HueOffset2       = 2.0 / 3.0
	HueOffset3       = 1.0 / 6.0
	HueOffset4       = 1.0 / 2.0
	HueOffset5       = 2.0 / 3.0

	// LSL model thresholds and values
	LSL_Threshold  = 1e-10
	LSL_HueRed     = 0.0
	LSL_HueYellow  = 60.0
	LSL_HueGreen   = 120.0
	LSL_HueCyan    = 180.0
	LSL_HueBlue    = 240.0
	LSL_HueMagenta = 300.0
	LSL_HueRange   = 60.0

	// Wavelength ranges for LSL model
	LSL_RedMin    = 700.0
	LSL_YellowMin = 570.0
	LSL_GreenMin  = 520.0
	LSL_CyanMin   = 495.0
	LSL_BlueMin   = 450.0
	LSL_VioletMin = 380.0

	// RGB8 model constants
	RGB8_Max = 255.0
)

// Test color values
const (
	AlmostWhiteR = 0.9999
	AlmostWhiteG = 0.9999
	AlmostWhiteB = 0.9999
	AlmostWhiteA = 1.0

	AlmostBlackR = 0.0001
	AlmostBlackG = 0.0001
	AlmostBlackB = 0.0001
	AlmostBlackA = 1.0

	Gray50R = 0.5
	Gray50G = 0.5
	Gray50B = 0.5
	Gray50A = 1.0

	TransparentWhiteR = 1.0
	TransparentWhiteG = 1.0
	TransparentWhiteB = 1.0
	TransparentWhiteA = 0.0

	SemiTransparentRedR = 1.0
	SemiTransparentRedG = 0.0
	SemiTransparentRedB = 0.0
	SemiTransparentRedA = 0.5

	YellowR = 1.0
	YellowG = 1.0
	YellowB = 0.0
	YellowA = 1.0

	CyanR = 0.0
	CyanG = 1.0
	CyanB = 1.0
	CyanA = 1.0

	MagentaR = 1.0
	MagentaG = 0.0
	MagentaB = 1.0
	MagentaA = 1.0

	PastelPinkR = 1.0
	PastelPinkG = 0.8
	PastelPinkB = 0.8
	PastelPinkA = 1.0

	PastelBlueR = 0.8
	PastelBlueG = 0.8
	PastelBlueB = 1.0
	PastelBlueA = 1.0

	PastelGreenR = 0.8
	PastelGreenG = 1.0
	PastelGreenB = 0.8
	PastelGreenA = 1.0
)

// Expected conversion values
const (
	AlmostWhiteLABL = 99.99
	AlmostWhiteLABA = 0.0
	AlmostWhiteLABB = 0.0

	AlmostBlackLABL = 0.01
	AlmostBlackLABA = 0.0
	AlmostBlackLABB = 0.0

	Gray50LABL = 53.39
	Gray50LABA = 0.0
	Gray50LABB = 0.0
	Gray50YUVY = 0.5
	Gray50YUVU = 0.0
	Gray50YUVV = 0.0
	Gray50YIQY = 0.5
	Gray50YIQI = 0.0
	Gray50YIQQ = 0.0

	TransparentWhiteLABL = 100.0
	TransparentWhiteLABA = 0.0
	TransparentWhiteLABB = 0.0

	SemiTransparentRedLABL = 53.24
	SemiTransparentRedLABA = 80.09
	SemiTransparentRedLABB = 67.20

	YellowHSLH = 60.0
	YellowHSLS = 1.0
	YellowHSLL = 0.5
	YellowHSBA = 1.0
	YellowHSBB = 1.0
	YellowLABL = 97.14
	YellowLABA = -21.55
	YellowLABB = 94.48
	YellowLSBW = 570.0
	YellowLSBS = 1.0
	YellowLSBB = 1.0
	YellowLSLL = 0.5

	CyanHSLH = 180.0
	CyanHSLS = 1.0
	CyanHSLL = 0.5
	CyanHSBA = 1.0
	CyanHSBB = 1.0
	CyanLABL = 91.11
	CyanLABA = -48.09
	CyanLABB = -14.13
	CyanLSBW = 495.0
	CyanLSBS = 1.0
	CyanLSBB = 1.0
	CyanLSLL = 0.5

	MagentaHSLH = 300.0
	MagentaHSLS = 1.0
	MagentaHSLL = 0.5
	MagentaHSBA = 1.0
	MagentaHSBB = 1.0
	MagentaLABL = 60.32
	MagentaLABA = 98.23
	MagentaLABB = -60.82
	MagentaLSBW = 380.0
	MagentaLSBS = 1.0
	MagentaLSBB = 1.0
	MagentaLSLL = 0.5

	PastelPinkHSLH = 0.0
	PastelPinkHSLS = 1.0
	PastelPinkHSLL = 0.9
	PastelPinkHSBA = 0.2
	PastelPinkHSBB = 1.0
	PastelPinkLABL = 86.41
	PastelPinkLABA = 18.00
	PastelPinkLABB = 6.87

	PastelBlueHSLH = 240.0
	PastelBlueHSLS = 1.0
	PastelBlueHSLL = 0.9
	PastelBlueHSBA = 0.2
	PastelBlueHSBB = 1.0
	PastelBlueLABL = 83.57
	PastelBlueLABA = 10.30
	PastelBlueLABB = -24.91

	PastelGreenHSLH = 120.0
	PastelGreenHSLS = 1.0
	PastelGreenHSLL = 0.9
	PastelGreenHSBA = 0.2
	PastelGreenHSBB = 1.0
	PastelGreenLABL = 95.46
	PastelGreenLABA = -25.58
	PastelGreenLABB = 19.18
)

// LAB model constants
const (
	LAB_Epsilon = 216.0 / 24389.0
	LAB_Kappa   = 24389.0 / 27.0
)

// D65 illuminant white point
const (
	LAB_RefX = 0.95047
	LAB_RefY = 1.00000
	LAB_RefZ = 1.08883
)

// RGB to XYZ conversion matrix
const (
	LAB_RGB_X1 = 0.4124564
	LAB_RGB_X2 = 0.3575761
	LAB_RGB_X3 = 0.1804375
	LAB_RGB_Y1 = 0.2126729
	LAB_RGB_Y2 = 0.7151522
	LAB_RGB_Y3 = 0.0721750
	LAB_RGB_Z1 = 0.0193339
	LAB_RGB_Z2 = 0.1191920
	LAB_RGB_Z3 = 0.9503041
)

// XYZ to RGB conversion matrix
const (
	LAB_XYZ_R1 = 3.2404542
	LAB_XYZ_R2 = -1.5371385
	LAB_XYZ_R3 = -0.4985314
	LAB_XYZ_G1 = -0.9692660
	LAB_XYZ_G2 = 1.8760108
	LAB_XYZ_G3 = 0.0415560
	LAB_XYZ_B1 = 0.0556434
	LAB_XYZ_B2 = -0.2040259
	LAB_XYZ_B3 = 1.0572252
)

// sRGB conversion constants
const (
	LAB_SRGB_LinearThreshold = 0.0031308
	LAB_SRGB_Gamma           = 2.4
	LAB_SRGB_LinearScale     = 12.92
	LAB_SRGB_GammaScale      = 1.055
	LAB_SRGB_GammaOffset     = 0.055
)

// HSL conversion constants
const (
	HSL_DegreesPerCircle = 360.0
	HSL_HueNormalization = 1.0 / 360.0
	HSL_HueOffset1       = 1.0 / 3.0
	HSL_HueOffset2       = 2.0 / 3.0
	HSL_HueOffset3       = 1.0 / 6.0
	HSL_HueOffset4       = 1.0 / 2.0
	HSL_HueOffset5       = 2.0 / 3.0
)

// LCH conversion constants
const (
	LCH_DegreesPerCircle   = 360.0
	LCH_RadiansPerDegree   = math.Pi / 180.0
	LCH_DegreesPerRadian   = 180.0 / math.Pi
	LCH_WhiteThreshold     = 1e-10
	LCH_NearWhiteThreshold = 0.01
	LCH_ChromaThreshold    = 1e-10
)

// HCL conversion constants
const (
	HCL_DegreesPerCircle   = 360.0
	HCL_RadiansPerDegree   = math.Pi / 180.0
	HCL_DegreesPerRadian   = 180.0 / math.Pi
	HCL_WhiteThreshold     = 1e-10
	HCL_NearWhiteThreshold = 0.1
	HCL_ChromaThreshold    = 1e-10
)

// CMY model constants
const (
	CMY_One           = 1.0
	CMY_PercentageMax = 100.0
)

// CMYK model constants
const (
	CMYK_One           = 1.0
	CMYK_PercentageMax = 100.0
)

// HEX model constants
const (
	HEX_Prefix      = "#"
	HEX_Max         = 255.0
	HEX_Base        = 16
	HEX_Bits        = 8
	HEX_ShortLength = 3
	HEX_LongLength  = 6
	HEX_ShortOffset = 1
	HEX_LongOffset  = 2
)

// Wavelength model constants
const (
	// Visible spectrum range in nanometers
	WavelengthMin = 380.0
	WavelengthMax = 750.0

	// Constants for wavelength to RGB conversion
	WavelengthGamma = 0.8

	// Wavelength ranges for RGB conversion (normalized to [0,1])
	WavelengthRange1 = 0.17 // Blue range
	WavelengthRange2 = 0.42 // Cyan range
	WavelengthRange3 = 0.58 // Green range
	WavelengthRange4 = 0.83 // Red range

	// RGB coefficients for wavelength conversion
	WavelengthBlueCoeff1  = 0.4
	WavelengthBlueCoeff2  = 0.6
	WavelengthGreenCoeff1 = 0.3
	WavelengthGreenCoeff2 = 0.7
	WavelengthRedCoeff1   = 0.3
	WavelengthRedCoeff2   = 0.7
	WavelengthRedCoeff3   = 0.6

	// Normalization factors for wavelength ranges
	WavelengthNorm1 = 0.17 // Blue normalization
	WavelengthNorm2 = 0.25 // Cyan normalization
	WavelengthNorm3 = 0.16 // Green normalization
	WavelengthNorm4 = 0.25 // Yellow normalization
	WavelengthNorm5 = 0.17 // Red normalization
)
