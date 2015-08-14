# MDT

Simple program to aid during meditation. Logs exact Hz plus label of mind occurrence.

**Note**: This is a small program I wrote for a friend of mine. It's quite 
unlikely that you will be able to use the program itself for any real purpose. 
Nevertheless, after taking permission, I decided to open source it, in case 
someone finds the code useful, maybe to write something with similar 
functionality. The code is a little messy so be careful.

Things that might interest you:

* Use of [termbox-go](https://github.com/nsf/termbox-go) to receive keys from 
  keyboard.
* Use of channels to synchronize timers based on key presses.
* Text based input that is enabled on mouse click.

## Usage

* Download the executable and place in a folder (it will produce logs and 
  config.json) or build your own:

    go get github.com/nstratos/mdt
    go build

* Click with the mouse on the configuration values (Like Mode, Offset etc) to
  change them.
* Press the spacebar to start the timer.
* After key capturing starts, record key presses (q, w, e, a, s or d).
* Either press spacebar to end the session or wait for the timer to finish.
* End the program anytime by pressing 'Esc'.
* View the log that was produced.

## Screenshots

![mdt changing configuration](/screenshots/mdt_input.png?raw=true "Changing configuration")

![mdt capturing key press](/screenshots/mdt_capturing.png?raw=true "Capturing key press")


## Specifications

The program should work on Windows and preferably it should be a standalone 
application and not a web one.

1.  Input total running time in two digit integer. Say '15' Have stop button 
    if stopped early.
2.  Input offset time in minutes, say '6' (integer 0-99)
    * Select which mode: A or B. Only one can be active.
    * Input for 3 digit integer for Base. Say '80'. Or '150'
3.  Input start hz in numbers, say '14,54' (2 decimals)
4.  Input end hz in numbers, say '15,38' (2 decimals).
5.  Press Start button. This starts the timer for total running time.

    Now the program will start and record keypresses for the q,w,e,a,s,d keys.

    Whenever one such key is pressed it will log the exact hz, the time plus 
the label associated with a key. (and maybe calculation Base hz too for linear 
base hz progression)

6.  'w' is pressed, log says: 15,05hz @ 80 base hz, on 04:30 Visual memory 

    It calculates the exact hz by math. It assumes a meditation runs linearly 
from start hz to end hz over the whole of its running time. So if starting at 
15 hz, going to 19 hz, in 20 minutes, if 'w' is pressed at 04:30 then the hz 
would be: 

	`hzPerSecond = (EndHz - StartHz) / (TotalTime - Offset) * 60`

	`secondsSinceOffset = currentSecs - (Offset * 60)`

	`currentHz = hzPerSecond * secondsSinceOffset + StartHz`

    Each key corresponds to a label.

    | key | label description    |
    | --- | -------------------- |
    | q   | visual memory        |
    | a   | visual imagination   |
    | w   | auditory memory      |
    | s   | auditory imagination |
    | e   | language voice       |
    | d   | language thought     |

    Use a newline after each logged occurence.

7.  Program stops after running time has passed, or when stopped by user 
    pressing Stop button.

8.  It will write a log file in .txt in its directory, preferably named: 
    S-E hz day date month time where S is start hz and E is end hz, so for 
    example file name: '15-19 hz wed 27 dec 22.09.txt'

    Put at the top of the text file the filename and the mode used (A or B)

