
package main

import (
    "io/ioutil"
    //"unsafe"
    "fmt"
    "encoding/binary"
    "bytes"
	"os"
)

type LuaHeaderType struct {
    W0, W1, W2 uint32
}

type LuaInstructionType uint32

//Holds part of a Function block for the top-level (main()). It is missing the source name
//Not included in this type due to variable length string.
type LuaTopLevelFuncType struct {
    LineDefined, LastLineDefined uint32
    NumUpVals, NumParameters byte // For the main chunk, the values of both fields are 0.
    IsVarg byte
    MaxStackSize byte
}

type LuaFuncType struct {
    SourceNameSize uint32//0 for non-main funcs
    LineDefined, LastLineDefined uint32
    NumUpVals, NumParameters byte
    IsVarg byte
    MaxStackSize byte
}

func extractOpcodeIndex(instr LuaInstructionType) LuaInstructionType {
    return instr & 0x3F
}

func main() {
    //------- READ ENTIRE FILE ---------------
	fileBytes, fileErr := ioutil.ReadFile(os.Args[1])
    if fileErr != nil { panic(fileErr) }

    var lhdr LuaHeaderType
    mybuf := bytes.NewBuffer(fileBytes[:])//fileBytes[4:]
    myerr := binary.Read(mybuf, binary.LittleEndian, &lhdr)
    if myerr != nil {
        fmt.Println("binary.Read failed:", myerr)
    }

    //-------- HEADER ---------------
    if lhdr.W0 != 0x61754C1B {
        panic ("Lua signature not found")
    }
    
    major := (lhdr.W1 >> 4)  & 0xF
    minor := (lhdr.W1 >> 0)  & 0xF
    fmt.Printf("Lua version %d.%d\n", major, minor)
    
    //-------- Original source filename
    var StringLength uint32
    binary.Read(mybuf, binary.LittleEndian, &StringLength)
    if StringLength != 0 {
        var srcFilename = make([]byte, StringLength)
        binary.Read(mybuf, binary.LittleEndian, &srcFilename)
        fmt.Printf("Src file is %s\n", srcFilename)
    }
    
    var mainFunc LuaTopLevelFuncType
    binary.Read(mybuf, binary.LittleEndian, &mainFunc)
    fmt.Printf("First line %d.\n", mainFunc.LineDefined)
    fmt.Printf("Last line %d.\n", mainFunc.LastLineDefined)
    fmt.Printf("Up vals %d.\n", mainFunc.NumUpVals)
    fmt.Printf("Parameters %d.\n", mainFunc.NumParameters)
    fmt.Printf("Varg %d.\n", mainFunc.IsVarg)
    fmt.Printf("Max stack size %d.\n", mainFunc.MaxStackSize)
    
    
    var codeSize uint32
    binary.Read(mybuf, binary.LittleEndian, &codeSize)
    fmt.Printf("Code size %d.\n", codeSize)
    
    var mainInstList = make([]LuaInstructionType, codeSize)
    binary.Read(mybuf, binary.LittleEndian, &mainInstList)
    
    var instNum uint32 = 0
    for ; instNum < codeSize; instNum++ {
    
     luaP_opnames := [...] string { "MOVE","LOADK","LOADBOOL","LOADNIL","GETUPVAL","GETGLOBAL",
  "MOVE",
  "LOADK",
  "LOADBOOL",
  "LOADNIL",
  "GETUPVAL",
  "GETGLOBAL",
  "GETTABLE",
  "SETGLOBAL",
  "SETUPVAL",
  "SETTABLE",
  "NEWTABLE",
  "SELF",
  "ADD",
  "SUB",
  "MUL",
  "DIV",
  "MOD",
  "POW",
  "UNM",
  "NOT",
  "LEN",
  "CONCAT",
  "JMP",
  "EQ",
  "LT",
  "LE",
  "TEST",
  "TESTSET",
  "CALL",
  "TAILCALL",
  "RETURN",
  "FORLOOP",
  "FORPREP",
  "TFORLOOP",
  "SETLIST",
  "CLOSE",
  "CLOSURE",
  "VARARG"}
    
        opcodeIdx := extractOpcodeIndex(mainInstList[instNum])
        fmt.Printf("Opcode is %s\n", luaP_opnames[ opcodeIdx ])
    }
    


}

//http://golangtutorials.blogspot.co.uk/2011/06/structs-in-go-instead-of-classes-in.html

