MDT
======

Hey man, thanks for wanting to program this for me. Idk, it might change a little, and I'll have to get names for the labels, and maybe the output should be in CSV so it can be imported into a spreadsheet, anyway this is what I've got so far.

Function
--------
Logs exact Hz plus label of mind occurence.

Program usage
-------------

1. Input total running time in two digit integer. Say '15' Have stop 		    button if stopped early.
2. Input offset time in minutes, say '6' (integer 0-99)
  * Select which mode: A or B. Only one can be active.
  * Input for 3 digit integer for Base. Say '80'. Or '150'
3. Input start hz in numbers, say '14,54' (2 decimals)
4. Input end hz in numbers, say '15,38' (2 decimals).
5. Press Start button. This starts the timer for total running time.

    Now the program will start and record keypresses for the q,w,e,a,s,d keys.

    Whenever one such key is pressed it will log the exact hz, the time   plus the label associated with a key. (and maybe calculation Base hz too for linear base hz progression)

6. 'w' is pressed, log says: 15,05hz @ 80 base hz, on 04:30 Visual memory 

    It calculates the exact hz by math. It assumes a meditation runs linearly from start hz to end hz over the whole of its running time. So if starting at 15 hz, going to 19 hz, in 20 minutes, if 'w' is pressed at 04:30 then the hz would be: 

    hz per second H = (19-15) / ((20-6) * 60) total seconds passed is from counter, say S.

    hz = S * H + start hz

    hz = S * (end hz - start hz) / ((total time - offset) * 60)

    the key corresponds to a label.

   key, label description:
   * q visual memory
   * a visual imagination
   * w auditory memory
   * s auditory imagination
   * e
   * d

   Please put a newline after each logged occurence.

7. Program stops after running time has passed, or when Stopped by user pressing Stop button.

8. It will write a log file in .txt in its directory, preferably named: S-E hz day date month time
where S is start hz and E is end hz, so for example file name: '15-19 hz wed 27 dec 22.09.txt'

  Please put at the top of the text file the filename and the mode used (A or B)

That's it. Thanks.

I'll think some more about the labels and the exact output text for the log file. Plus calculate linear Base hz as well. Maybe Comma Seperated Values for import in 'excel' open office?﻿