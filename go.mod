module github.com/golubaca/heif-converter

go 1.23.3

require (
	github.com/MaestroError/go-libheif v0.3.0
	github.com/strukturag/libheif v1.16.2 // indirect
)

require github.com/golubaca/goheif v0.0.0-20241112162159-463fdc52e1a0

require github.com/jdeng/goheif v0.0.0-00010101000000-000000000000 // indirect

replace github.com/jdeng/goheif => github.com/golubaca/goheif v0.0.0-20241112232426-fb239e86afe7
