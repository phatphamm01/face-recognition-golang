#Use macbook

## 1. Install golang

```sh
brew install go
```

## 2. Install pkg-config, opencv, dlib, and jpeg

```sh
brew install opencv
brew install dlib
brew install jpeg
```

## 3. Config environment

```sh
export CGO_CXXFLAGS="--std=c++11"

export CGO_CPPFLAGS="-I/opt/homebrew/Cellar/opencv/4.8.1_5/include -I/opt/homebrew/Cellar/dlib/19.24.2/include -I/opt/homebrew/Cellar/jpeg/9e/include"

export CGO_LDFLAGS="-L/opt/homebrew/Cellar/jpeg/9e/lib -ljpeg -L/opt/homebrew/Cellar/dlib/19.24.2/lib -ldlib -L/opt/homebrew/Cellar/opencv/4.8.1_5/lib -lopencv_stitching -lopencv_superres -lopencv_videostab -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dpm -lopencv_face -lopencv_photo -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_surface_matching -lopencv_tracking -lopencv_datasets -lopencv_dnn -lopencv_plot -lopencv_xfeatures2d -lopencv_shape -lopencv_video -lopencv_ml -lopencv_ximgproc -lopencv_calib3d -lopencv_features2d -lopencv_highgui -lopencv_videoio -lopencv_flann -lopencv_xobjdetect -lopencv_imgcodecs -lopencv_objdetect -lopencv_xphoto -lopencv_imgproc -lopencv_core"
```

## 4. Install go package

```sh
go install
```

## 5. Run

```sh
go run main.go
```
