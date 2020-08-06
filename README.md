# addlicense
This project adds the given license to all source file in a given directory.
The tool will try to avoid source files that already has a header. This
project can be used as a CLI as well as a library imported to a go script.

## Inspiration
The initial inspriation is from this [repo](https://github.com/google/addlicense) from Google. While
the google addlicense project is good, there are some additional requirements I needed.
1. I want the code to be a importable library in addition to just a tool. In
this way, I can import this to write my own customized scripts. 
2. Template is an overkill for most people. Most of the time, I just want to
attach a license to all my files. If I am making changes to the license, such
as updating the year, I will just update a single license file. So I
simplified the codebase.

## Usage
