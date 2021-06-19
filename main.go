package main

import (
    "bufio"
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/png"
    "io/ioutil"
    "log"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "golang.org/x/image/math/fixed"
)

func main() {
    transparent := flag.Bool("t", false, "背景色を透明にします")
    flag.Parse()

    // ログファイルの準備
    NowTime := time.Now()
    LogFileName := NowTime.Format("20060102_150405_") + ".txt"
    logfile, err := os.OpenFile(LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }
    defer logfile.Close()

    // テキストファイルの読み込み
    f, err := os.Open("emoji.txt")
    if err != nil {
        log.Fatal(err)
        return
    }

    s := bufio.NewScanner(f)
    cnt := 0

    // テキストファイルを行ごとに読み込んでループ   
    for s.Scan() {
        Message := s.Text()
        OutputFileName := "output/" +  strconv.Itoa(cnt) + ".png"

        // カンマ区切りなら後半をファイル名にする
        if strings.Contains(s.Text(), ",") {
            slice := strings.Split(s.Text(), ",")
            Message = slice[0]
            OutputFileName = "output/" +  slice[1] + ".png"
        }
 
        // 画像サイズを決める
        img := image.NewRGBA(image.Rect(0, 0, 128, 128))

        // 画像背景色を決める
        for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
            for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
                if *transparent {
                    img.Set(x, y, color.RGBA{0, 0, 0, 0})
                } else {
                    img.Set(x, y, color.RGBA{200, 200, 200, 200})
                }
            }
        }

        // フォントファイルの読み込みとパース
        ftBin, err := ioutil.ReadFile("fonts/NotoSansCJKjp-Bold.otf")
        if err != nil {
            log.Fatalf("failed to load font: %s", err.Error())
        }
        ft, err := opentype.Parse(ftBin)
        if err != nil {
            log.Fatalf("failed to parse font: %s", err.Error())
        }

        // 文字数を判定する
        if len(Message) / len([]rune(Message)) != 3 {
            fmt.Fprintln(logfile, Message + " ERROR[ごめんなさい！全角文字を使ってね！]")
            continue
        }

        // 文字数による描画設定(フォントサイズ、行数、描画座標)
        var CharSize float64 = 64
        CharSize = 64
        Row1Xvalue,Row1Yvalue,Row2Xvalue,Row2Yvalue := 0,0,0,0
        Row1Text,Row2Text := "", ""
        if len([]rune(Message)) == 2 {
            CharSize = 64
            Row1Xvalue = 0 * 64
            Row1Yvalue = 85 * 64
            Row1Text = Message
        } else if len([]rune(Message)) == 3 {
            CharSize = 42
            Row1Xvalue = 0 * 64
            Row1Yvalue = 75 * 64
            Row1Text = Message
        } else if len([]rune(Message)) == 4 {
            CharSize = 64
            Row1Xvalue = 0 * 64
            Row1Yvalue = 60 * 64
            Row2Xvalue = 0 * 64
            Row2Yvalue = 120 * 64
            Row1Text = Message[0:6]
            Row2Text = Message[6:12]
        } else if len([]rune(Message)) == 6 {
            CharSize = 42
            Row1Xvalue = 0 * 64
            Row1Yvalue = 55 * 64
            Row2Xvalue = 0 * 64
            Row2Yvalue = 105 * 64
            Row1Text = Message[0:9]
            Row2Text = Message[9:18]
        } else if len([]rune(Message)) == 8 {
            CharSize = 32
            Row1Xvalue = 0 * 64
            Row1Yvalue = 50 * 64
            Row2Xvalue = 0 * 64
            Row2Yvalue = 100 * 64
            Row1Text = Message[0:12]
            Row2Text = Message[12:24]
        } else{
            fmt.Fprintln(logfile, Message + " ERROR[ごめんなさい！全角2,3,4,6,8文字以外は対応してません！]")
            continue
        }

        // 文字色の候補
        color := []*color.RGBA{
            {255, 0, 0, 255},
            {255, 219, 0, 255},
            {73, 255, 0, 255},
            {0, 255, 146, 255},
            {0, 146, 255, 255},
            {73, 0, 255, 255},
            {255, 0, 219, 255},
        }

        // 文字色をランダムにあてがう
        rand.Seed(time.Now().UnixNano())
        RandomNumber := rand.Intn(7)

        opt := opentype.FaceOptions{
            Size:    CharSize,
            DPI:     72,
            Hinting: font.HintingNone,
        }
        face, err := opentype.NewFace(ft, &opt)

        // 1行目の描画
        d := &font.Drawer{
            Dst: img,
            Src: image.NewUniform(color[RandomNumber]),
            Face: face,
            Dot: fixed.Point26_6{fixed.Int26_6(Row1Xvalue), fixed.Int26_6(Row1Yvalue)},
        }
        d.DrawString(Row1Text)

        // 2行目の描画
        if len([]rune(Message)) >= 4 {
            d2 := &font.Drawer{
                Dst: img,
                Src: image.NewUniform(color[RandomNumber]),
                Face: face,
                Dot: fixed.Point26_6{fixed.Int26_6(Row2Xvalue), fixed.Int26_6(Row2Yvalue)},
            }
            d2.DrawString(Row2Text)
        }

        // ファイル出力
        file, err := os.Create(OutputFileName)
        if err != nil {
            panic(err.Error())
        }
        defer file.Close()
        fmt.Fprintln(logfile, Message + " " + OutputFileName)
        
        if err := png.Encode(file, img); err != nil {
            panic(err.Error())
        }
        
        cnt ++
    }

    if s.Err() != nil {
        // non-EOF error.
        log.Fatal(s.Err())
    }

}