# Sound of Sort

https://github.com/user-attachments/assets/0fe81ea9-1ae5-420b-8c6c-aabd3ddf87ea

This is a Go-based terminal application that visualizes and sonifies a whole bunch of sorting algorithms, inspired by [Sound of Sorting](https://panthema.net/2013/sound-of-sorting/). 

---

## What it Does
- **Real-time Visualization:** Renders sorting algorithms as a bar graph directly in your terminal.
- **Sonification:** Every time an array element is accessed or modified, it plays a tone (pitch corresponds to the element's value)
- **Interactive:** Change algorithms, speed, volume, array size and reshuffle on the fly.
- **Image Mode:** Pipe ASCII/ANSI art via stdin with `-img`
- **Algorithms:** Includes a loads of classic (and not-so-classic) sorting algorithms (you can add your own too)

---

## Available Algorithms

- Quick Sort
- Bubble Sort
- Selection Sort
- Insertion Sort
- Merge Sort
- Heap Sort
- Shell Sort
- Cocktail Shaker Sort
- Gnome Sort
- Pancake Sort
- Radix Sort (LSD)
- Timsort
- Bitonic Sort
- Bogo Sort

---

## Installation

### Prerequisites
* Go (version 1.21+ recommended)

```bash
go install github.com/sahaj-b/sound-of-sort@latest
# binary will be installed in $GOPATH/bin or $GOBIN
```

## Building from source

```bash
# Clone the repository
git clone https://github.com/sahaj-b/sound-of-sort.git
cd sound-of-sort

go build
```

-----

## Usage

```bash
./sound-of-sort
```

### Command-Line Flags

You can customize the startup state. If you don't, it uses sane defaults.

| Flag      | Description                                         | Default   |
| :-------- | :-------------------------------------------------- | :-------- |
| `-sort`   | Initial sorting algorithm to use                    | `quick`   |
| `-size`   | Initial array size (ignored in image mode)          | `100`     |
| `-delay`  | Initial delay between operations (ms)               | `5`       |
| `-volume` | Initial volume (0.0 to 1.0)                         | `0.1`     |
| `-fps`    | Rendering frames per second                         | `60`      |
| `-list`   | List all available sorting algorithms and exit      | `false`   |
| `-img`    | Enable image mode (pipe ASCII art via stdin)        | `false`   |
| `-horiz`  | Horizontal image mode (sort rows instead of columns) | `false`   |
| `-help`   | Show this help message and exit                     | `false`   |

Image mode: when `-img` is set the program reads an ANSI/ASCII image from stdin, and uses its width as the array length. The `-size` flag is ignored in this mode.

**Example:** Start with the Bogo Sort on a tiny array

```bash
./sound-of-sort -sort bogo -size 8
```

-----

## Image Mode (ASCII / ANSI Art Sorting)

Pipe any ASCII/ANSI colored art into this beast with `-img`.  
By default, individual columns of the image are shuffled and sorted (the image visually reassembles).  

Basic usage:

```bash
pixcii -i path/to/image.jpg -c | sound-of-sort -img
# or
chafa -f symbols path/to/img.jpg -s 50x50 | sound-of-sort -img -horiz
# or
ascii-image-converter path/to/image.jpg | sound-of-sort -img
# or any other ASCII art generator
```

-----

## Controls


| Key(s)                 | Action                          |
| :-------------------   | :----------------------------   |
| `q` or `Ctrl+C`        | **Quit** the application        |
| `←` / `→` or `p` / `n` | Cycle through algorithms        |
| `↑` / `↓`              | Increase / Decrease volume      |
| `w` / `s`              | Increase / Decrease delay       |
| `a` / `d`              | Decrease / Increase array size  |
| `r`                    | **Reshuffle** the current array |


---

## Adding a New Sorting Algorithm

### Step 1: Create the Algorithm File
First, create a new file in the `/algos` directory. eg: `algos/ur_cool_sort.go`.

### Step 2: Write Your Sort Function

Inside your new file, write your sorting function. It **must** follow this signature:

```go
package algos

import "context"

func urCoolSort(ctx context.Context, arr ArrObj) {
    // Your sorting logic goes here
    n := arr.Len()

    // Example:
    for i := range {
        if arr.Get(ctx, i) > arr.Get(ctx, i+1) {
            arr.Swap(ctx, i, i+1)
        }
    }
}
```

**Pay attention, this is the important part:**
- Use the `arr` object for all array operations. You've got `arr.Get(ctx, index)`, `arr.Set(ctx, index, value)`, `arr.Swap(ctx, i, j)`, and `arr.Len()`.
- **DO NOT** try to access a raw slice. The `ArrObj` interface is what hooks your algorithm into the visualizer and sound engine.

### Step 3: Register Your Algorithm

Now, tell the application your new sort exists. Open `algos/main.go` and find the `Sorts` slice. Add your algorithm to the list.

```go
// in algos/main.go

var Sorts = []struct {
    Name string
    Arg  string
    Fun  sortFunc
}{
    {"Quick Sort", "quick", quickSort},
    {"Bubble Sort", "bubble", bubbleSort},
    // ... all the other sorts
    {"Your Cool Sort", "cool", urCoolSort},
    //                    │      ^^^^^^^ your main sorting function
    //                    ╰ this is for command-line arg
}
```

### Step 4: Rebuild and Run

```bash
go build
./sound-of-sort --sort cool
```
