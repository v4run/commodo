/**
@author: varun<varunvasan91@gmail.com>
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	dir     string               // Root directory for file server
	port    string               // Port on which file server should run
	version bool                 // Display version
	help    bool                 // Display help
	icons   = map[string]string{ // Icons for different file types
		"css":       "AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAunABAAAAAAAAAAAAunABILpwAUm6cAF+unAB2LhtAMm5bgBvunABQbtxARgAAAAAAAAAALpwAQAAAAAAAAAAAAAAAAC6cAGCunAB37pwAfW6cAH/unAB/7xzBP/Ulx3/xYEO/7htAP+2awDyunAB07pwAV4AAAAAAAAAALpwAQAAAAAAunAB/7pwAf+6cAH/um8A/7lpAP+7bwD/26Ib/+CoIP/iryz/15wg/8aCDv+5bgDGAAAAAAAAAAAAAAAAAAAAALpwAf+6cAH/um8G/7+HNv/Js5P/09TY//v58//w2Jv/47NE/96mIf/boyT/t2sA5wAAAAAAAAAAAAAAALpwARS6cAH/um0D/9HX4P/R1+H/0djg/8/Jv//9/P7////////////68t//3aQg/7htAP8AAAAAunABAAAAAAC6cAEzunAA/7t1Ff/R2OL/y76p/7hnAP+5aAD/2p0Q/+rFcP/+/fv//vz3/+CpIv+6bwH/unABAgAAAAAAAAAAunABTrpwAP+7dhX/zMGu/8SeZf+6bgD/vHQE/92nKf/fqSj//frx///////hrCb/vHME/7pwARcAAAAAAAAAALpwAXK6cAH/u3II/7pzF/+7dhn/vHof/758If/frjn/4bA7//357///////4a0u/794B/+6cAE4AAAAAAAAAAC6cAGOuWwA/8KXWf/Q0tX/z87M/8/Ozf/Rz8//+/v5///9+////////////+KvNv/CfAn/unABTQAAAAAAAAAAunABrrlrAP/EnWb/z83L/9DU2f/Pz8//0tLS//v48f/36sr/9OCy//bnxP/hsTz/w30K/7pwAXoAAAAAAAAAALpwAce6cAH/unAC/7ltBP/DmF3/0NXb/9LS1P/8/Pv//v35//LdrP/dpR3/36gn/8WBDv+6cAGTAAAAAAAAAAC6cAHfuW0A/8OZXf/Dm2D/w5da/8OYW//Nv6v//P39/////////////////+jCaP/GgQn/unABsQAAAAC6cAEBunAB57pwBP/R1t7/z8/O/8/Pz//Pz9D/0tLS//z8+////v3///79///////v1pn/x4IF/7pwAccAAAAAunABEbpwAe66cAP/vX8p/71+Jv+9fif/vn4n/7+BKf/iskX/47VH/+O1R//ktUj/4rA9/8uKEv+5bwDnu3ECB7pwAR26cAHzunAB/7pwAf+6cAH/unAB/7pwAf+8dAT/3agn/9+qKf/fqin/36op/9+qKf/Rkxn/t20A7rtxARK6cAEtunAB+rpwAf66cAH+unAB/rpwAf66cAH+unAB/rtxAv67cQL+u3EC/rtxAv67cQL+u3EC/rpwAfa6cAEi/n8AAMAHAADAAwAAwAMAAMADAADAAwAAwAMAAMADAACAAwAAgAMAAIABAACAAQAAgAEAAIABAACAAQAAgAEAAA==",
		"html":      "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAATdJREFUeNpifOqtlsDAwDCfgQwgvfUWIxOQfsBAAWCiQO8BEMECxBdgIvypVQzc/vF4db2tjGX4efkUnM8C9McHYDiAOf++fgbTv+9dZ/gPZaODfwjxDzAXYICvmxYyfNuznpAXLiIbAPKGwa/LJ4FUDgOXcxADs5gMVl2fl01G4bMgOwcG2HTNwBgd/H31FNmABxhe+H3vBjiQsAHh9sUMf18+RRZCMeAgEDv8+/oJJYThcc3NR3w6kNpyk0F83j4GdqgXmMWlGThdAsHsP69wu+ADIpo+AQNQGuxkjDBA8gIw+h+gxwIYvM4LYOCwcAG7AESDDPxxYi/Dj+N7GH5h8R4LNlu+blwIxnjABfQwADlnAZEZC6R5IozDiC4LTNYGoBgBYn8oDQqfDdCY2gBK+sjqAQIMACfPddVHcozkAAAAAElFTkSuQmCC",
		"js":        "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAB8lBMVEX65FEoKi4nKS734VCmmUN+dzz541D75FHx3E8qLC8rLS/75VEmKS4sLi/v208mKC7z3k/541Hz3lDbyEz131D85VH+51Hv2k8gIy324FD951Hy3U8uLy8wMTAxMjDHtkgpKy/r1k6jl0OVikCVi0Csn0T44lAyMzCflEL85lHn007031D03lAvMDAlKC6ckEHz3U+GfT4fIi1XUza2qEb34lC0pkVFRDN6czvy3E9lXzjGtkjiz01IRzMhJC1/eDzYxku5qkbx3Vg+PjKypEXRv0rz4Wzy4Wt6czzz4nCkl0Pn000qKy+jl0Lx3VcvMC/95lHCs0eXjEHHt0gkJi5aVzdbVzcjJi7q1k5HRjN/dz2Rh0AtLi9BQDI4ODFHRTM5OTH/6VLBsUe2p0Wek0K9rkclJy5cWDf54lDPvkr24U/GtUg0NTDOvUotLy/75lHw2012cDutoESglELt2E/Ds0iSiEDo1E765FAiJS3y3mDhzk0uMC9tZzpWUzY9PTLn1E4rLC/75VC3qEbWxEtMSjSqnUT/8VPMu0ny4GgsLS+OhD8nKi51bjtYVTZiXTijlkP/61KYjUHcyUxJRzSlmEM/PjL24VDTwkqLgj785VBYVDYkJy7k0E3u2U9pZDnx3l/y3mHx3lvx3lzw207w208rnsfNAAAAdUlEQVQY02NwW4gCqhgWLUUBS2gqkGkiJJA7SZgDLsDKyN5dI8dlhiSQGNcqJi+CJMDtIVenhaJl+kyeLhWoAAOfPSM7m4BUGH8TWCBJIrueN98ovWNaDxtYQHCKnkLL0sCYNBnhJQxOIC0CmmxAUk3TdOkSAP2UgcdJHqHTAAAAAElFTkSuQmCC",
		"music":     "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAqFJREFUeNp0Us9PE1EQnre77W5Zty1SWlD5YY1RxJig0YSTHIheTQz1P/BC5OBFE6PEK4kJRogelQux4sEE46UXLhghGIWQaGJYEVqbbbW4LVu23d3nvN1CCwkv+bL75s18880PQimFRAoOHlr7JhEziPcI46BTchCAg0NOKrU0sbKi9mtaIek4Tg5NrxFDiKZGP+Ewgnj82B38jCwvr/XzPDcUiYRvxmLhRCQSAo7j5vBtYI/AMOi+YA51+XwCbW0NJePx9plCofSgUCjeZWT4fG1w8NLori9hPeh7Xj1AQKiqfn1MCEkEg3JPNNpsRKPh2UDAv9uPbRbLeuASXHhmQC7PNSqgjrNNTPMnu/YiEkg2pChyT0tLELq7293kb6/XmtgcsqGri4cqCmEwTYBKRQZCZPa8ihjFROd0vXReVTP3slkbNM2ul3B1qgg+KQA7Ozz8zlJAEyqghDno/77sK0+SToIoypQQIGv3xf1TCAQIdHYQMMoAjuPZgsEzsLnxvaamFdU1oTqWhNbHSKkz4NjWGDbvoiBwvKIQe2+V4AiIGEj5MGh5bwVYLCqsE+h/temOeFvb0YgPy6iClrVUTOk6lEu6p05yoPOEBYYp1QgaFimfTctKswL6loHQQRDEp7zgQycbzLIBflFynX2CA2HJws4JUKvAm8LWH+3R8sf54urSYjqX2RyxbXsiPfsEfr0bc52QaMCyKgu2VbWoXaGEVm2eWHUFjm2PVz6/GTd5Hky/COn5qfpOxM6C1NU3fbr3VFvsuAxGyYSN9Zy6/WMB4NYNb4zEq5dnKhH+Gti/+xC8PfOt9/IVhV20TAZ21pYebr4afoGx+cYxMgLJbbuHwG6JpdXU5GJJHxYEoUjWP70sz01+QHMLa99/AQYAohQjWPbGXSYAAAAASUVORK5CYII=",
		"pdf":       "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAALEgAACxIB0t1+/AAAABx0RVh0U29mdHdhcmUAQWRvYmUgRmlyZXdvcmtzIENTNui8sowAAAAWdEVYdENyZWF0aW9uIFRpbWUAMTAvMzAvMDZ8tDgbAAACAUlEQVQ4jX2TsUtbURTGf+/lmegjTQOhMZg2QzM51Kil/4CjEBwECaKjg6NLFhcHd6eC4B+QpbgoOEoGkXRoHVtF+pQWqhaDFvTdPHNPh9z3eKa2Bz7O5b3vnu8799xrXafTkgBsA8uATAYZHES32+ggoAtog65ILwN2FwgMOoACfMBeXye1v8/D2Bgq/C9CRyTiB4ATiPQUQ2UT7ugo13NzWLkcyiiKUY3WItjKqIbqoQP/6Iihep07z8MPvxsHyqADWF9cV+LKYnJqfJxXe3t8LRZB60g1nrUITueJzQJ0Ly7QWiOlEsrzkFgbOlbI+jg0JPSFAIV6HUkkcCcm+La0RHBz88jBX2fgG9wDHdclu7jI960t7s7OKG5soPN5HjIZHtJpdDaLL4ICHCXyyD62TWFhgfbhIYlymd+eR2l2lkqrxeX2NurqimeVCp9rNbQI9r1R1bkcamAAGRmhsLyMnc/zvFrl5vSUZqXCz50dEsPDXDab+O02vgg+YH1IpQTg3e4ux2trFOfncctlWjMzkauw75e1GoXpaX4dHHCyuYkAjjIkdXuLOznJi2qV5tQUyowpXuSk0eC40YgOUQOOb0g6meT1ygqfVle5Pj9/NFIdW0eTMPus98mk9N+BiGxITyFsy1Exm9Jf4F8IXdN7P2+BEeBN/4X6T9wDl8CPP7gaJqQXV/CHAAAAAElFTkSuQmCC",
		"image":     "R0lGODlhEAAQAPcAAEKU50Kl91JCQmNCKWNjY2sxIWtjtWtzpXOEtXOl53O173uMxnuUxnul3oQAhIRzpYSczoStzoSt3oxaMYyEtYyUtZSEpZSUxpSlzpTG95yl1qWElKWUlKW956XO960xKa1rIa2Ura21/63W/7UxELVaMbVzQrWEa7WclLXW/72Me73O/8aclMbW/8bn/84xGM6Uc9ZCGNacc9alc9be59bv/95KId6UWt61hN7n7+dSIeeUUufOhOfn5++EOe/Ge+/v9+/3//e9Y/fGc/9rMf9zMf+EOf+EUv+MQv+UQv+UY/+tUv+1Wv/GY//We//enP/ne//nhP///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////yH+EkR1c3RpbiBGcmllc2VuaGFobgAh+QQBAAAOACwAAAAAEAAQAAAIwgAlCBw4kICDgwglSFnIUIoEBQYROlDYcGGAAAAiHqTYMEiLFiIQaFSohAgRJEuaBKlRY8VIKUdKFEnC5MmQJRhcJpTyoQCJFzF06LABQefGIDsmgPBhRIgRHxGMTgQyw8SAGzucOLkhQSqEHixUCDghg8cPGA2kMsgR4oGBDTyiQEHRwIPGBTkuULBggQOOHxYS2EWIgEaNESM8eMjAeEQHjYVruJicojLixwgPaNisAcOFzxcqUNDogIDp06hNIwwIADs=",
		"xml":       "R0lGODlhEAAQAPcAAG6Ov3KQv3KRv3aTvXqVu3uVvH6XulWAyFmBxliCxV2ExF6ExGGIw2KJw2WKwmeLwWmMwWyOwP9mAP8A/4KavIWdvoigwIuiwoyiwo+lxJGnxpKoxpWryJesyZmtypCu1p2xzaO20am71K/A2LPD2tLi+tTj+tbk+tfl+9nm+9vn+9zo+97p++Dq/OHs/OPt/OXu/Obv/Ojw/erx/evy/e3z/e/0/fD2/vL3/vT4/vX5/vf6/vn7/////xLuYAAAQAAAAAAAABLuqBLuaAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQBLurJDuGAAAAAAAAgAAAQAC7BLuYAAC7JEZcCPwGCPv+AAAAAAAAAHskAAAjxLsBJDuGBLswNS5TACpEQCpERLs1NS47rV3GACpEQAAABLs7NS5srV3GNcVPRQCgBQCQBLtDNdNrxQCgMAEttdN4xQCgBLtFAAAAJEFyCBguBLt4JEFURQHqJEFbRLuOAAAABLtPAAAAJEFyCKQSBLuCJEFURQHSBLtWAAAAJEFyCKQSBLuJJEFURQHSJEFbRLuaAAABAAAAOaERAAAAgAABAAAMAAAACBgwNSLsf3wAAAAMAAABBQAABLrmJD7bAAAIAAAACKQUBLuOAAAAAAAIADwqgAAIAAAAAAAAJDnvJDVhhLuCJD7bJD7cZDVhpDnvBQAABLt5JDnyBLujJDuGJD7eAH//wAABBLtaAAAABLujJDuGJEFcP///5EFbZEJvBQAAAAAACKQUBLuSJEJkiKQUAAAABLunN3tDt3tIGKmyAACWGKm1AAAAAAAAAAAAAAAAAAAABLuaBLu7BRyIBRyIBLuoOb8I8OlLsYaoBLu2MLCzQAABMLC4xR9cBRyIAAAAxRwGsXS4BRwABLu1BLupP///xLvQMNclMEgcP///8LC40SV1RRyIGMboGMboEUEtRQAABR9cJLOwAAAAAAAAOqG1OqG1OqG1OqG1AACpBLvJN1sdBLvLJLOwJLOwObgowAACeaCsAAABCH5BAEAABMALAAAAAAQABAAAAinACcIJCEiBIgOGi5UMGBAoMMJI3js0JEDxw0bBno0nCCho8ePHjVy7CGBJMmOJktqLMmyh8mUDA2gnJmS5AsXMmuGVOmihUyQQCW0YEHAoYcaNGbIiAHj5tAVAxxySPqhqtUPK1QEcLhBKVOnLLKmiOAww9KrVlOgeOAQw9eeYVWoPdHAoYWmaKueMLHAYQWwYlHsLYHAIYUCBQYIAADBAQMFCQ4cCAgAOw==",
		"txt":       "R0lGODlhEAAQAKIFAIWFhYSEhAAAAMbGxv///////wAAAAAAACH5BAEAAAUALAAAAAAQABAAAAM9WBXcrnCRSWcQUY7NSQAYFFSVYIYapxIDOpJVK7JqJysvPN1pPbAuHYU38m2AMyESR/MtJUqiUeU6Wa+FBAA7",
		"video":     "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABGdBTUEAAK/INwWK6QAAABl0RVh0U29mdHdhcmUAQWRvYmUgSW1hZ2VSZWFkeXHJZTwAAAJYSURBVHjaYvz//z8DJQAggBiINICJhYVFR0JCIsfR0XF7UVHxM0ZGRhWQXoAAYsGjSYaHh8fU3t7ekYOD09bQ0jP1dWDadKkfoZ//5h+MDMzqwDV3AEIIEaQKUDTQBqEhYSEjMTFxR2UlVUcs7JyjLdu3cLm6urOcXqNYZlyxbdZGPjPHLx4rkD//79PQFU/wio9xdAAIEN8PX13SQiImYOBGIaGtoMXV1t/9XUt5u3Lj2+IMH9/cDFR8F4nvYvA8QQGADEhKSz0hLyxrv3r39969fvy8/vTwyK9fPy8zMjK94Obm+v3nzx9WJiZmNmZmRiagPCMHBwfDhw8f3nz9+nU/QACBw+Dk0e3ra2tjIODg1kdHBwNV6xYYRQTE81w5MgxBk5OdgZdXX2GNWtWMcTHJzBs3bqVQUVFmeHgwYMPgVoVAAKICWTAuXPn1l6+fPm1goISw48fvhBGl69esugoKAINICb4efPXwyamjoML168ZlBT02D4+xes5xZIL0AAwaOxsrL6BJD9/9On7/+/f/3//Pnn/+/ffv9/8ePv2CxHz8gYiD+9+9//puams0C6QUIIORoZAKZDPQjGIMANHbAlsDEfv/+zCNOWYQHyCA4Ab8+wdS9A+sAAaQDUCIgXz9Hy4GEEBIBvyD2v4HyVGMsAhDMZSREaECIIBYkOMUpnVAEwAMoCJiRHuAoAAghvAzMzCxs7OwcDG9oOgARAXMLKCSIAAghuwd+/uiz9+/OD8+fMn3tzFxQEjNaf/+7cuXUExAcIIEZKszNAgAEA/zob40WaW00AAAAASUVORK5CYII=",
		"directory": "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAABX1BMVEX39/dra2tsbGxtbW1ubm57e3t+fn5/f3+ua1Cua1G2cVO4dFS7eFW9e1a6fVe/f1e7f1i9g1m/hFnBgljDhVnFiFnGiVrHjFrDiF7HjV/IjFrIjFvJjlvKkFzLkVzMk13NlF3PlV3Oll7Pl17Rml/SnV/Pmmbfl2HTnWDUn2DUoGHVoWHWo2TZqGXZqWfRoGrToGvXqW/0rWn0t2z5tWv+uW3+vG7hrXT6zHitra2vr6+wsLCysrLJrZnLtZ/ivYbhvpDXvKrZvqn90oP33IP52oH/1orlxJPjw5nu15z45Y385pn56p3dxavjxqDjyKr/76j7763//aP05rry47z88LD89L3//77MzMzU1NTW1tbd3d399sf/+Mf48cv9+Mz9+c727dL9+dn//9v+/t3l5eXn5+fo6Ojp6en16+T48Ob9+uX//OX//uX+/up9W1D8+PX4+Pj5+fn///j///+gwhWcAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAAHdElNRQfVAxgNNRYCF8YoAAAAr0lEQVR42mNgIAwcdDT1fdMQfD/37LxEA7touIBhXEpKSoiajLiYhCNYnVZ8AhQEC1imAgVUwzz1tDVUlORlpcUEbYFqVAJ1s6C6czL5rNIYlH20i5NiwMA/l5c1kkHBSyM5FAI8YnmYIhjkXFXCvSHAPoCbMYJBxkIpyBkCzN24gAKSxoouZqYgYKJuxAHU4iQkKiUmIszPy8PNxZnPFsWQbs3GwgwFLOw2GRh+BQApPielSWe5LQAAAABJRU5ErkJggg==",
		"home":      "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAetJREFUeNpi/P//PwM2sFNepA9IpUC5c9wfvinCpo4RmwHb5CCaXXfP5AXxd7umfwYZciIhiwFq6JympiaIgSADkPFGGeE+IP70/dyM/xuUpP6BMIgNEpvmaPPz69ev/4uLiz/B1LMg27xWGmKzx8Ii3i1BTf9V4hIYQeIgtuf8It6tSZP+X6irZvj48QtCE8yk5ZLCfUD86f283P/L5CX/namt+v/7928wBrFBYiC5lUqy/5uN9H7C9IGJheJCfUD86XVD6P+FshL/TtZU/v/16xcKBonNkxb7D1IzT0L4J0gP2IA5YkJ9QPzpebLt/zkyEv+OVlf8//nzJ1YMkgOpAasF6gHpZZguIvTpsbfW/+lSEv8OVpb///79O14MUgNSC9YD1MvswckpeuPJF53fhoaMftNmMDMQANL2Dgyn9u//fevUrW//gB5gyXnzFhSfRe3t7Z/+/fvHBlOYnZ2FonHq1Glw9jd7h5+VO3bygdjwaAQGFAPQABRNtrZ2YPrw4UMociC1MAA3ABhdDOip8s+fPyjRjawWwwCQYnQX/P4NsgmcllDkkA3G64KfP38RdAE4M33Z0/b/w/yNKJprWDUY1NTUwexbt24ytPy+gSIvkOjPwONSxQg2YK6k8H8GMkDy87eMAAEGAFwSczT/jIkHAAAAAElFTkSuQmCC",
		"file":      "R0lGODlhEAAQANUvAH5+fuXl6NHR1f39/f///+Hh5PPy9O7u8Orp7NjX3N3c4Pb2987N083N0tXU2Onp7Nzc4PLy9NXU2djY3OHh4/r6+/f29+/u8Pb3+NHQ1c7N0vr5+/r6+vb2+NnX3Pn5++np7fn6+u7t8fPz9OHg5Orq7NXT2NnX3dTU2e7u8eTk5NTU2N3d4N3c39TT2dbr9QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAEAAC8ALAAAAAAQABAAAAaLwNcLQCwahcghYclcEpPDgXQqBagASUDoU+F0N5ujEGDBLMyLdMfSwY4jEYPcMIrL3cND6nC48A8ifngACCUgCAgPiogIgwGPkJEBgyQFFAWYl5gFgwoKLRAsnhAQnoMJCScJE6isHhODDg4SDisSLg4oJhKDGQK/wMAZgwwMDcUaDQ0axoNGz0UvQQA7",
	}
)

const (
	VERSION = "1.0"
	NAME    = "Commodo"
	COMMAND = "commodo"
	// DIR_IMAGE  = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAABX1BMVEX39/dra2tsbGxtbW1ubm57e3t+fn5/f3+ua1Cua1G2cVO4dFS7eFW9e1a6fVe/f1e7f1i9g1m/hFnBgljDhVnFiFnGiVrHjFrDiF7HjV/IjFrIjFvJjlvKkFzLkVzMk13NlF3PlV3Oll7Pl17Rml/SnV/Pmmbfl2HTnWDUn2DUoGHVoWHWo2TZqGXZqWfRoGrToGvXqW/0rWn0t2z5tWv+uW3+vG7hrXT6zHitra2vr6+wsLCysrLJrZnLtZ/ivYbhvpDXvKrZvqn90oP33IP52oH/1orlxJPjw5nu15z45Y385pn56p3dxavjxqDjyKr/76j7763//aP05rry47z88LD89L3//77MzMzU1NTW1tbd3d399sf/+Mf48cv9+Mz9+c727dL9+dn//9v+/t3l5eXn5+fo6Ojp6en16+T48Ob9+uX//OX//uX+/up9W1D8+PX4+Pj5+fn///j///+gwhWcAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAAHdElNRQfVAxgNNRYCF8YoAAAAr0lEQVR42mNgIAwcdDT1fdMQfD/37LxEA7touIBhXEpKSoiajLiYhCNYnVZ8AhQEC1imAgVUwzz1tDVUlORlpcUEbYFqVAJ1s6C6czL5rNIYlH20i5NiwMA/l5c1kkHBSyM5FAI8YnmYIhjkXFXCvSHAPoCbMYJBxkIpyBkCzN24gAKSxoouZqYgYKJuxAHU4iQkKiUmIszPy8PNxZnPFsWQbs3GwgwFLOw2GRh+BQApPielSWe5LQAAAABJRU5ErkJggg=="
	// FILE_IMAGE = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAAwFBMVEXz9fv4+fyDkbGDkbCHk66DkrGCkbCDkrDu8vry9fv5+vz4+fuGmr6Mn8KGlK6Llqvw9PuTpsaarMuis8+ouNOsvNXe6Pjd5/ff6Pfk7Pnn7vnt8vrz9vvc5/ff6fje6Pfh6vjq8Pn2+PuQmqjm7vnw9fvy9vuWnaT1+Pv3+fuVnaSboaGhpJynqJinqJmtq5Wyr5KyrpK7tIu3sY63sY+/tojgyI/UsmjUsmnavXzVsmn///8AAAAAAAAAAAAAAAAYiYSvAAAAPHRSTlP//////////////////////////////////////////////////////////////////////////////wC7iOunAAAAAWJLR0Q/PmMwdQAAAAlwSFlzAAAASAAAAEgARslrPgAAAI1JREFUGNNdyukSgiAYhWGxlHbbTYHKtMW0bKOy7v++Ug5OUy/Dn+c7xvsvo/hS94WW6ikruMSkHROSv6SG22446HZW9/xxBZybfMMNXnQCZOvJeFSO+hngwFiDqY6ANNCLIAUkka+K/ASw9/TC2wLCZVUIcObqPltMHUBPCLd8rqgDTIvaFqU2NWuAnz7Z5yiNWyg7MwAAAABJRU5ErkJggg=="
	// HOME_IMAGE = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAetJREFUeNpi/P//PwM2sFNepA9IpUC5c9wfvinCpo4RmwHb5CCaXXfP5AXxd7umfwYZciIhiwFq6JympiaIgSADkPFGGeE+IP70/dyM/xuUpP6BMIgNEpvmaPPz69ev/4uLiz/B1LMg27xWGmKzx8Ii3i1BTf9V4hIYQeIgtuf8It6tSZP+X6irZvj48QtCE8yk5ZLCfUD86f283P/L5CX/namt+v/7928wBrFBYiC5lUqy/5uN9H7C9IGJheJCfUD86XVD6P+FshL/TtZU/v/16xcKBonNkxb7D1IzT0L4J0gP2IA5YkJ9QPzpebLt/zkyEv+OVlf8//nzJ1YMkgOpAasF6gHpZZguIvTpsbfW/+lSEv8OVpb///79O14MUgNSC9YD1MvswckpeuPJF53fhoaMftNmMDMQANL2Dgyn9u//fevUrW//gB5gyXnzFhSfRe3t7Z/+/fvHBlOYnZ2FonHq1Glw9jd7h5+VO3bygdjwaAQGFAPQABRNtrZ2YPrw4UMociC1MAA3ABhdDOip8s+fPyjRjawWwwCQYnQX/P4NsgmcllDkkA3G64KfP38RdAE4M33Z0/b/w/yNKJprWDUY1NTUwexbt24ytPy+gSIvkOjPwONSxQg2YK6k8H8GMkDy87eMAAEGAFwSczT/jIkHAAAAAElFTkSuQmCC"
)

func init() {
	flag.StringVar(&dir, "d", getHomeDir(), "The root directory for the file server.")
	flag.StringVar(&dir, "directory", getHomeDir(), "The root directory for the file server.")
	flag.StringVar(&port, "p", "4545", "The port on which the file server should run.")
	flag.StringVar(&port, "port", "4545", "The port on which the file server should run.")
	flag.BoolVar(&version, "v", false, "Prints the version number.")
	flag.BoolVar(&version, "version", false, "Prints the version number.")
	flag.BoolVar(&help, "h", false, "Prints the version number.")
	flag.BoolVar(&help, "help", false, "Prints the version number.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", COMMAND)
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\t-d, -directory  Directory   The root directory for the file server.\n")
		fmt.Fprintf(os.Stderr, "\t-p, -port       Port        The port on which the file server should run.\n")
		fmt.Fprintf(os.Stderr, "\t-v, -version    Version     Prints the version number.\n")
		fmt.Fprintf(os.Stderr, "\t-h, -help       Help        Show this help.\n")
	}

}

func showVersion() {
	fmt.Println("\n", NAME, VERSION)
	fmt.Println("This is a free software and comes with NO warranty.\n")
}

type fileServerHandler struct {
	root http.FileSystem
}

func htmlDocumentBegin() string {
	return "<html>\n" +
		"<head>\n" +
		"<title>" + NAME + "</title>\n" +
		"<style>" +
		"body{margin:0; padding-top: 10px}" +
		".contents{padding-left:10px}" +
		".entry{padding-left: 18px }" +
		"a:link{color: #2c2c2c; text-decoration: none}" +
		"a:visited{color: #2c2c2c; text-decoration: none}" +
		"a:hover{color: black; text-decoration: underline}" +
		"a:active{color: #2c2c2c; text-decoration: none}" +
		"</style>\n" +
		"</head>\n" +
		"<body><div class = 'contents'>\n"
}

func htmlDocumentEnd() string {
	return "</body>\n</html>"
}

func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (handler *fileServerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	requestedPath := request.URL.Path
	if !strings.HasPrefix(requestedPath, "/") {
		requestedPath = "/" + requestedPath
		request.URL.Path = requestedPath
	}

	f, openErr := handler.root.Open(path.Clean(requestedPath))
	if openErr != nil {
		http.Error(writer, "404 Page not found.", http.StatusNotFound)
		return
	}

	dir, statErr := f.Stat()
	if statErr != nil {
		http.Error(writer, "404 Page not found.", http.StatusNotFound)
		return
	}
	defer f.Close()

	url := request.URL.Path
	if dir.IsDir() {
		if url[len(url)-1] != '/' { // show contents of directory
			localRedirect(writer, request, path.Base(url)+"/")
			return
		}
	} else {
		if url[len(url)-1] == '/' { // show file even if a / is found after filename
			localRedirect(writer, request, "../"+path.Base(url))
			return
		}
	}

	contentType := mime.TypeByExtension(filepath.Ext(dir.Name()))
	if !dir.IsDir() {
		writer.Header().Set("Content-Type", contentType)
	} else {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	if dir.IsDir() {
		fmt.Fprintf(writer, htmlDocumentBegin())
		var folders string
		var files string
		for {
			dirs, err := f.Readdir(100)
			if err != nil || len(dirs) == 0 {
				break
			}
			for _, d := range dirs {
				name := d.Name()
				if strings.HasPrefix(name, ".") {
					continue
				}
				if d.IsDir() {
					folders += fmt.Sprintf("<span class = 'folder'><img src=data:image/png;base64,"+icons["directory"]+" /> <a href=\"%s\">%s</a></span><br>\n", name, name)
				} else {
					// if image, found := icons[filepath.Ext(d.Name())]; found {
					// 	files += fmt.Sprintf("<span class = 'entsry' ><img src=data:image/png;base64,"+image+" /> <a href=\"%s\">%s</a></span><br>\n", name, name)
					// } else {
					files += fmt.Sprintf("<span class = 'enstry' ><img src=data:image/png;base64,"+icons["file"]+" /> <a href=\"%s\">%s</a></span><br>\n", name, name)
					// }
				}
			}
		}
		fmt.Fprintf(writer, folders)
		fmt.Fprintf(writer, files)
	} else {
		if request.Method != "HEAD" {
			io.CopyN(writer, f, dir.Size())
		}
		return
	}

	fmt.Fprintf(writer, "</div><hr style='border: 0; background-color: #5c5c5c; height:0.5px'>\n")
	fmt.Fprintf(writer, "\n<a href=\"/\" style = 'border-right: 1px dashed #9c9c9c; padding: 8.5px; margin-right: 10px;'> <img src=data:image/png;base64,"+icons["home"]+" /></a><span style='font-family: \"Times New Roman\"; color: #2c2c2c; font-style:italic; font-size:14;'>Served using Commodo v%s</span>\n", VERSION)
	fmt.Fprintf(writer, "<hr style='border: 0; background-color: #5c5c5c; height:0.5px'>\n")

	fmt.Fprintf(writer, htmlDocumentEnd())
}

func main() {

	go handleExit() // goroutine to handle ctrl + c
	flag.Parse()    // parse command line arguments

	if version {
		showVersion()
		os.Exit(0)
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	_, fileErr := os.Stat(dir)
	if fileErr != nil { // Check if path exists
		fmt.Println("Invalid Path `", dir, "`. Please specify a valid path.")
		os.Exit(1)
	}

	startServer() // start the file server
}

func startServer() {
	fmt.Printf("Starting %s with root %s.\nPress ctrl + c to exit.", strings.Title(NAME), dir)
	http.Handle("/", &fileServerHandler{http.Dir(dir)})
	conErr := http.ListenAndServe(":"+port, nil)
	if conErr != nil {
		fmt.Println(conErr)
		os.Exit(1)
	}
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user info.")
		os.Exit(1)
	}
	return usr.HomeDir

}

func handleExit() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	for sig := range signalChannel {
		if sig == syscall.SIGINT {
			fmt.Printf("\n %s File server stopped.\n", NAME)
			os.Exit(0)
		}
	}
}
