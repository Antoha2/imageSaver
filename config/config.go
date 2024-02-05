package config

const ImgPath = "D:\\code\\images\\buffer\\"

var ImgFormats = [...]string{"png", "jpeg", "jpg", "bmp", "gif", "tiff"}

type ContextKey string

const ImgMap ContextKey = "ImgMap"

const HostAddr = ":8180"
