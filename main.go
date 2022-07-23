package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dslipak/pdf"
	"github.com/signintech/gopdf"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	orderPosition       = 666
	fileExistErrorMsg   = "unedited csv file exists."
	fileUnexistErrorMsg = "unedited csv file exists."
	pdfPath             = "./inputData/入力データ.pdf"
	csvPath             = "./outputData/中間データ_未編集.csv"
	editedCsvPath       = "./outputData/中間データ.csv"
	boxType             = "/MediaBox"
	resultPdfDir        = "./outputData/pdf/"
	fontPath            = "./src/font/ipag.ttf"
	costFontSize        = 18
	detailFontSize      = 9
	stampPath           = "./src/img/stamp.jpg"
)

const (
	processingUnitPrice = 3500
	toolText            = "梱包一式 段ボール"
	costCoefficient     = 2.8
)

type OrderExtracted []struct {
	PageNumber  int
	Name        string
	Code        string
	Type        string
	DueDate     string
	Amount      string
	OrderNumber string
}

type DataSet struct {
	X       float64
	Y       float64
	Content string
}

type PreparedData struct {
	PageNumber  int
	OrderNumber string
	Year        DataSet
	Month       DataSet
	Date        DataSet

	m_Cost DataSet
	p_Cost DataSet
	t_Cost DataSet

	m_Unit  DataSet
	p_Unit  DataSet
	t_Whole DataSet

	t_Depl DataSet

	UnitCost  DataSet
	TotalCost DataSet

	m_Breakdown DataSet
	p_Breakdown DataSet
	t_Breakdown DataSet
}

func (Orders OrderExtracted) convertOrders() [][]string {
	ordersSlice := make([][]string, len(Orders)+1)
	orderInfo := reflect.TypeOf(Orders[0])

	ordersSlice[0] = []string{
		orderInfo.Field(0).Name,
		orderInfo.Field(1).Name,
		orderInfo.Field(2).Name,
		orderInfo.Field(3).Name,
		orderInfo.Field(4).Name,
		orderInfo.Field(5).Name,
		orderInfo.Field(6).Name,
		"材料単価",
		"材料費",
		"材料名",
		"加工単価",
		"加工費",
		"加工費内訳のとこ",
		"要具費",
		"要具費内訳のとこ",
	}

	for i, order := range Orders {
		ordersSlice[i+1] = []string{
			strconv.Itoa(order.PageNumber),
			order.Name,
			order.Code,
			order.Type,
			order.DueDate,
			order.Amount,
			order.OrderNumber,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		}
	}
	return ordersSlice
}

func prepareData(csvFile *os.File) ([]PreparedData, error) {
	preparedData := []PreparedData{}

	csvReader := csv.NewReader(transform.NewReader(csvFile, japanese.ShiftJIS.NewDecoder()))
	for i := 0; ; i++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		preparedData = append(preparedData, PreparedData{})
		if i == 0 {
			continue
		}

		pd := &preparedData[i]
		pd.PageNumber, err = strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		pd.OrderNumber = record[6]

		now := time.Now()
		pd.Year.Content = strconv.Itoa(now.Year())
		pd.Month.Content = strconv.Itoa(int(now.Month()))
		pd.Date.Content = strconv.Itoa(now.Day())

		est_M_cost, err := strconv.ParseFloat(record[8], 32)
		if err != nil {
			return nil, err
		}
		pd.m_Cost.Content = strconv.FormatFloat(est_M_cost*costCoefficient, 'f', 2, 32)

		est_P_cost, err := strconv.ParseFloat(record[11], 32)
		if err != nil {
			return nil, err
		}
		pd.p_Cost.Content = strconv.FormatFloat(est_P_cost*costCoefficient, 'f', 2, 32)

		est_T_cost, err := strconv.ParseFloat(record[13], 32)
		if err != nil {
			return nil, err
		}
		pd.t_Cost.Content = strconv.FormatFloat(est_T_cost*costCoefficient, 'f', 2, 32)

		mUnit, err := strconv.ParseFloat(record[7], 32)
		if err != nil {
			return nil, err
		}
		pd.m_Unit.Content = fmt.Sprintf("%.2f", mUnit)

		pUnit, err := strconv.ParseFloat(record[10], 32)
		if err != nil {
			return nil, err
		}
		pd.p_Unit.Content = fmt.Sprintf("%d", int(pUnit))

		amount, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			return nil, err
		}
		pd.t_Depl.Content = record[5]

		mCost, err := strconv.ParseFloat(pd.m_Cost.Content, 32)
		if err != nil {
			return nil, err
		}
		pCost, err := strconv.ParseFloat(pd.p_Cost.Content, 32)
		if err != nil {
			return nil, err
		}
		tCost, err := strconv.ParseFloat(pd.t_Cost.Content, 32)
		if err != nil {
			return nil, err
		}
		pd.t_Whole.Content = fmt.Sprintf("%.2f", amount*tCost)

		pd.UnitCost.Content = strconv.FormatFloat(mCost+pCost+tCost, 'f', 2, 32)

		unitCost, err := strconv.ParseFloat(pd.UnitCost.Content, 32)
		if err != nil {
			return nil, err
		}
		totalCost := amount * unitCost
		pd.TotalCost.Content = fmt.Sprintf("%.2f", math.Round(totalCost))

		mConsumption := mCost / mUnit * amount
		pd.m_Breakdown.Content = fmt.Sprintf("%s  %.2fkg", string(record[9]), mConsumption)

		pConsumption := pCost / (pUnit / 60) * amount
		pd.p_Breakdown.Content = fmt.Sprintf("%s成形機  %.2f分", string(record[12]), pConsumption)

		pd.t_Breakdown.Content = fmt.Sprintf("%s", record[14])

		// ---------- Prepare Positions ----------
		pd.Year.X = 433.0
		pd.Year.Y = 28.0

		pd.Month.X = 485.0
		pd.Month.Y = 28.0

		pd.Date.X = 520.0
		pd.Date.Y = 28.0

		pd.m_Cost.X = 370.0
		pd.m_Cost.Y = 345.0

		pd.p_Cost.X = 370.0
		pd.p_Cost.Y = 458.0

		pd.t_Cost.X = 370.0
		pd.t_Cost.Y = 570.0

		pd.m_Unit.X = 290.0
		pd.m_Unit.Y = 345.0

		pd.p_Unit.X = 290.0
		pd.p_Unit.Y = 458.0

		pd.t_Whole.X = 290.0
		pd.t_Whole.Y = 570.0

		pd.t_Depl.X = 241.0
		pd.t_Depl.Y = 570.0

		pd.UnitCost.X = 316.0
		pd.UnitCost.Y = 205.0

		pd.TotalCost.X = 428.0
		pd.TotalCost.Y = 205.0

		pd.m_Breakdown.X = 80.0
		pd.m_Breakdown.Y = 345.0

		pd.p_Breakdown.X = 80.0
		pd.p_Breakdown.Y = 458.0

		pd.t_Breakdown.X = 80.0
		pd.t_Breakdown.Y = 570.0
	}

	return preparedData[1:], nil
}

func main() {

	if _, err := os.Stat(editedCsvPath); err == nil {
		csvFile, err := os.Open(editedCsvPath)
		if err != nil {
			panic(err)
		}
		defer csvFile.Close()

		preparedData, err := prepareData(csvFile)
		if err != nil {
			panic(err)
		}

		fmt.Println(preparedData)
		fmt.Println("Processed data above into pdf files.")

		if err = createPdfs(preparedData); err != nil {
			panic(err)
		}

		return
	}

	if _, err := os.Stat(csvPath); err == nil {

		panic(fileExistErrorMsg)

	} else {

		content, err := readPdf(pdfPath)
		if err != nil {
			panic(err)
		}

		err = createCsv(content)
		if err != nil {
			panic(err)
		}
		fmt.Println("Csv-file is created.")
		// fmt.Println(content)
		return
	}

}

func createPdfs(preparedData []PreparedData) error {

	for _, item := range preparedData {
		pdf := gopdf.GoPdf{}
		pageSize := *gopdf.PageSizeA4
		pdf.Start(gopdf.Config{PageSize: pageSize})
		pdf.AddPage()
		templatePage := pdf.ImportPage(pdfPath, item.PageNumber, boxType)
		pdf.UseImportedTemplate(templatePage, 0, 0, pageSize.W, pageSize.H)

		// drawGrid := func(pdf *gopdf.GoPdf, page *gopdf.Rect) {
		// 	ww := 10.0
		// 	for i := 1; i < int(page.H/ww); i++ {
		// 		if i%10 == 0 {
		// 			pdf.SetLineWidth(0.8)
		// 			pdf.SetStrokeColor(50, 50, 100)
		// 		} else {
		// 			pdf.SetLineWidth(0.3)
		// 			pdf.SetStrokeColor(100, 100, 130)
		// 		}
		// 		x, y := float64(i)*ww, float64(i)*ww
		// 		pdf.Line(x, 0, x, page.H)
		// 		pdf.Line(0, y, page.W, y)
		// 	}
		// }
		//
		// // pdf.Line(0, 0, pageSize.W, pageSize.H)
		// drawGrid(&pdf, &pageSize)

		if err := pdf.AddTTFFont("ipag", fontPath); err != nil {
			return err
		}

		drawText := func(pdf *gopdf.GoPdf, x float64, y float64, s string) {
			pdf.SetX(x)
			pdf.SetY(y)
			pdf.Cell(nil, s)
		}

		// pdf.SetFont("ipag", "", costFontSize)
		// drawText(&pdf, 290, 207, "1")
		//
		// pdf.SetFont("ipag", "", detailFontSize)
		// drawText(&pdf, 80, 345, "ちゃんなな　ちゃんなな　ちゃんなな　ななななな")
		pdf.SetFont("ipag", "", detailFontSize)
		drawText(&pdf, item.Year.X, item.Year.Y, item.Year.Content)
		drawText(&pdf, item.Month.X, item.Month.Y, item.Month.Content)
		drawText(&pdf, item.Date.X, item.Date.Y, item.Date.Content)
		drawText(&pdf, item.m_Cost.X, item.m_Cost.Y, item.m_Cost.Content)
		drawText(&pdf, item.p_Cost.X, item.p_Cost.Y, item.p_Cost.Content)
		drawText(&pdf, item.t_Cost.X, item.t_Cost.Y, item.t_Cost.Content)
		drawText(&pdf, item.m_Unit.X, item.m_Unit.Y, fmt.Sprintf("%s /kg", item.m_Unit.Content))
		drawText(&pdf, item.p_Unit.X, item.p_Unit.Y, fmt.Sprintf("%s /時間", item.p_Unit.Content))
		drawText(&pdf, item.t_Whole.X, item.t_Whole.Y, fmt.Sprintf("%s", item.t_Whole.Content))

		drawText(&pdf, item.t_Depl.X, item.t_Depl.Y, item.t_Depl.Content)

		drawText(&pdf, 370.0, 660.0, item.UnitCost.Content)

		// drawText(&pdf, item.UnitCost.X, item.UnitCost.Y, item.UnitCost.Content)
		// drawText(&pdf, item.TotalCost.X, item.TotalCost.Y, item.TotalCost.Content)

		drawText(&pdf, item.m_Breakdown.X, item.m_Breakdown.Y, item.m_Breakdown.Content)
		drawText(&pdf, item.p_Breakdown.X, item.p_Breakdown.Y, item.p_Breakdown.Content)
		drawText(&pdf, item.t_Breakdown.X, item.t_Breakdown.Y, item.t_Breakdown.Content)

		intAndFraction := strings.Split(item.UnitCost.Content, ".")
		for i := 0; i <= 1; i++ {
			drawText(&pdf, item.UnitCost.X+float64(i*10), item.UnitCost.Y, string(intAndFraction[1][i]))
		}
		pdf.SetFont("ipag", "", costFontSize)
		for i := len(intAndFraction[0]) - 1; i >= 0; i-- {
			drawText(&pdf, -3.5+item.UnitCost.X+float64((i-len(intAndFraction[0])))*10.9, 2+item.UnitCost.Y, string(intAndFraction[0][i]))
		}

		pdf.SetFont("ipag", "", detailFontSize)
		intAndFraction = strings.Split(item.TotalCost.Content, ".")
		for i := 0; i <= 1; i++ {
			drawText(&pdf, item.TotalCost.X+float64(i*10), item.TotalCost.Y, string(intAndFraction[1][i]))
		}
		pdf.SetFont("ipag", "", costFontSize)
		for i := len(intAndFraction[0]) - 1; i >= 0; i-- {
			drawText(&pdf, -3.5+item.TotalCost.X+float64((i-len(intAndFraction[0])))*10.9, 2+item.TotalCost.Y, string(intAndFraction[0][i]))
		}

		if err := pdf.Image(stampPath, 500, 60, nil); err != nil {
			return err
		}

		if err := pdf.WritePdf(resultPdfDir + item.OrderNumber + ".pdf"); err != nil {
			return err
		}
	}
	return nil

}

func readPdf(path string) (OrderExtracted, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return nil, err
	}

	totalPage := r.NumPage()
	Orders := make(OrderExtracted, totalPage)

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		Orders[pageIndex-1].PageNumber = pageIndex
		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			if row.Position == orderPosition {
				// fmt.Println(">>>> row: ", row.Position)
				for i, word := range row.Content {
					switch i {
					case 0:
						Orders[pageIndex-1].Name = word.S
					case 1:
						Orders[pageIndex-1].Code = word.S
					case 2:
						Orders[pageIndex-1].Type = word.S
					case 3:
						Orders[pageIndex-1].DueDate = word.S
					case 5:
						Orders[pageIndex-1].Amount = word.S
					case 6:
						Orders[pageIndex-1].OrderNumber = word.S
					}
					// fmt.Println(word)
				}
			}
		}
	}
	// fmt.Println(Orders)
	return Orders, nil
}

func createCsv(Orders OrderExtracted) error {
	ordersSlice := Orders.convertOrders()

	for i, row := range ordersSlice {
		if i == 0 {
			continue
		}
		row[10] = strconv.Itoa(processingUnitPrice)
		row[14] = toolText
	}

	csvFile, err := os.Create(csvPath)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(transform.NewWriter(csvFile, japanese.ShiftJIS.NewEncoder()))
	for _, record := range ordersSlice {
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return nil
}
