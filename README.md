# PdfQuotationWriter
PDF-writer for streamlining a *horribly boring* paperwork.  
めんどうな事務仕事を効率化するためのpdfライターです。


## Why： 開発動機は？
Oftentimes, there exists a frigging boring tasks in manufacturing industries. Would you put up with it or to destroy them all? No need to ask, huh.<br>
これは多々あることですが、製造業にはとんでもなく退屈な仕事あったりします。　そんなもの我慢なんてしたくないので自動化してます。　　


## How： どのような対応？
There are excellent modules created by magnificent team or engineers. Why don't you utilize it?  
以下、素晴らしい技術者・チームの皆さんが開発したリソースを利用してます。  
- https://github.com/signintech/gopdf
- https://pkg.go.dev/github.com/dslipak/pdf
and so on.


## What: 何を製作したのか？
This program needs a specific-type-of-PDF file. So it's not intended to be used in general usage.<br>
このプログラムを動作させるためには、指定のフォーマットのpdfファイルが必要です。（そのままだと汎用性なし）　<br>
However, in automobile parts manufacturing industry, Hiroshima in Japan, You'll find a better way to manage your *boring paperwork*. That's what I wanted to solve.<br>
ここ日本、広島の自動車業界において、退屈な事務仕事をもっと効率的に処理できる方法を模索してます。　これこそが解決したかった課題です。個人的に<br>


Usage:　使用方法
1. Put *the* pdf in "./inputeData" directory.<br>
"./inputeData" ディレクトリに所定のpdfデータを保存してください。
2. Run the program with either CLI or GUIs.<br>
CLIでもGUIでもお好きな方で、プログラムを実行してください。
3. Check out results in "./outputaData/pdf"<br>
"./outputData/pdf"に出力されるファイルを確認してください。
then, upload each result to a corresponding request on customers' EDI.<br>
あとはEDIにアップロードするなりなんなり、お好きなように。


Step.1:  
Place input-pdf in "./inputData" directory.<br>
"./inpudaData"ディレクトリにファイル保存

![image_1](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_1.png "image_1")


Step.2:  
Run the compiled executable file or command "go run ." in the root of this project directroy.<br>
実行可能ファイルかコマンドで"go run ."を実行。　

![image_2](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_2.png "image_2")  

Then, an intermediate file gets created.<br>
中間ファイルが作成されます。

![image_3](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_3.png "image_3")  

Edit the intermediate file, typing in appropriate information like unit prices and remarks.<br>
中間ファイルに適切な情報を入力してくdささい。各種単価や備考情報などなど。

![image_4](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_4.png "image_4")  

Run the program again, the compiled one or just the command.<br>
再度、同様にプログラムを実行してください。

![image_5](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_5.png "image_5")  


Step.3:  
Pdf-files gets generated from the input-file in the directory "./outputData/pdf/", with respect to its P/O number.<br>
"./outputData/pdf/"ディレクトリにそれぞれの注文番号に対応したpdfファイルが出力されます。

![image_6](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_6.png "image_6")  

So like<br>
こんな感じで

![image_7](https://github.com/Kazuki-Maehara/PdfQuotationWriter/blob/image/image_7.png "image_7")  

After that, you can do whatever you like with these output-files at you disposal.<br>
このあとは、お好きなように　どうぞ。

