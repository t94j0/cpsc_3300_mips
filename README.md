# CPSC 3300: Instruction Set Simulator for MIPS-like Computer 

You should program a behavioral simulation of a simple MIPS-like instruction set.
The MIPS instruction set is covered in your textbook. We will use a subset of the instructions and the same instruction formats. However, there are multiple simplifications:

* 32-bit memory words with word addressability;
* a limited memory size of 1024 words;
* branch offsets and targets are not shifted before use;
* no jump or branch delay slots (i.e., jumps and branches have immediate effect);
* the program starts execution at address zero; and,
* no traps/exceptions/interrupts.

The instructions you should implement are:

```
opcode op/funct  action
  ------ --------  ------
  addu   0x00/0x21 r[rd]<-r[rs]+r[rt]
  addiu  0x09/n.a. r[rt]<-r[rs]+sign_ext(immed)
  and    0x00/0x24 r[rd]<-r[rs]&r[rt]
  beq    0x04/n.a. if(r[rs]==r[rt]) pc<-pc+sign_ext(immed)
  bgtz   0x07/n.a. if(signed(r[rs])>0) pc<-pc+sign_ext(immed)
  blez   0x06/n.a. if(signed(r[rs])<=0) pc<-pc+sign_ext(immed)
  bne    0x05/n.a. if(r[rs]!=r[rt]) pc<-pc+sign_ext(immed)
  hlt    all zero (which would be a nop in MIPS)
  j      0x02/n.a. pc<-target
  jal    0x03/n.a. r31<-updated_pc; pc<-target
  jalr   0x00/0x09 r[rd]<-updated_pc; pc<-r[rs]
  jr     0x00/0x08 pc<-r[rs]
  lui    0x0f/n.a. r[rt]<-immed<<16
  lw     0x23/n.a. r[rt]<-mem[r[rs]+sign_ext(immed)]
  mul    0x1c/0x02 r[rd]<-r[rs]*r[rt]
  nor    0x00/0x27 r[rd]<-~(r[rs]|r[rt])
  or     0x00/0x25 r[rd]<-r[rs]|r[rt]
  sll    0x00/0x00 r[rd]<-r[rt]<<shamt
  slti   0x0a/n.a. r[rt]<-(signed(r[rs])<sign_ext(immed))?1:0
  sra    0x00/0x03 r[rd]<-r[rt]>>shamt (sign bit duplicated)
  srl    0x00/0x02 r[rd]<-r[rt]>>shamt with zero fill
  subu   0x00/0x23 r[rd]<-r[rs]-r[rt]
  sw     0x2b/n.a. mem[r[rs]+sign_ext(immed)]<-r[rt]
  xor    0x00/0x26 r[rd]<-r[rs]^r[rt]
  xori   0x0e/n.a. r[rt]<-r[rs]^zero_ext(immed)
```

The instruction classifications are:
* alu ops: addu, addiu, and, lui, mul, nor, or, sll, slti, sra, srl, subu, xor, xori
* load: lw
* store: sw
* jumps: j, jr
* jump-and-links: jal, jalr
* branches: beq, bgtz, blez, bne
* halt: hlt

The instructions and data are read as hex values from stdin (e.g., using scanf() format specifier %x in C). The contents of memory are echoed as they are read in before the simulation begins; the contents are also displayed when a halt instruction is executed so that the changes to memory words caused by store instructions can be verified.

There are 32 registers, each 32 bits in size. Note that r0=0, as in regular MIPS.

A simple program to find the difference of two numbers, c = a - b, is shown below; a jump over an initial data area is used.
```
start: j    main
a:     0x22
b:     0x23
c:     0x0
main:  lw   r1,a
       lw   r2,b
       subu r3, r1, r2
       sw   r3, c
       hlt
```

The input file corresponding to this simple program is:

```
08000004
22
23
0
8c010001
8c020002
00221823
ac030003
0
```

Running the simulator with this input file (e.g., ./a.out < in1) results in this output:

```
contents of memory
addr value
000: 08000004
001: 00000022
002: 00000023
003: 00000000
004: 8c010001
005: 8c020002
006: 00221823
007: ac030003
008: 00000000
```

behavioral simulation of simple MIPS-like machine
  (all values are shown in hexadecimal)

```
pc   result of instruction at that location
000: j     - jump to 0x00000004
004: lw    - register r[1] now contains 0x00000022
005: lw    - register r[2] now contains 0x00000023
006: subu  - register r[3] now contains 0xffffffff
007: sw    - register r[3] value stored in memory
008: hlt

contents of memory
addr value
000: 08000004
001: 00000022
002: 00000023
003: ffffffff
004: 8c010001
005: 8c020002
006: 00221823
007: ac030003
008: 00000000

instruction class counts (omits hlt instruction)
  alu ops             1
  loads/stores        3
  jumps/branches      1
total                 5

memory access counts (omits hlt instruction)
  inst. fetches       5
  loads               2
  stores              1
total                 8

transfer of control counts
  jumps               1
  jump-and-links      0
  taken branches      0
  untaken branches    0
total                 1
```

As another example, here is a simple loop:

```
n = 5;
sum = 0;
for( i = 1; i <= n; i++ ){
  sum = sum + i;
}
```

The register allocation is: r1 = sum, r2 = i, r3 = n, and r4 = temp. (r0 is always 0.)

Here is MIPS-like assembly code for the loop:

```
start: addiu r3, r0, 5  // n = 5
       addu  r1, r0, r0 // sum = 0
       addiu r2, r0, 1  // i = 1
loop:  addu  r1, r1, r2 // sum = sum + i
       addiu r2, r2, 1  // i = i + 1
       subu  r4, r2, r3 // temp = i - n
       blez  r4, loop   // branch if temp <= 0
       hlt
 ```

After hand-assembly, the input is:

```
24030005
00000821
24020001
00220821
24420001
00432023
1880fffc
00000000
```

Running the simulator:
```
contents of memory
addr value
000: 24030005
001: 00000821
002: 24020001
003: 00220821
004: 24420001
005: 00432023
006: 1880fffc
007: 00000000
```

behavioral simulation of simple MIPS-like machine
  (all values are shown in hexadecimal)

pc   result of instruction at that location
```
000: addiu - register r[3] now contains 0x00000005
001: addu  - register r[1] now contains 0x00000000
002: addiu - register r[2] now contains 0x00000001
003: addu  - register r[1] now contains 0x00000001
004: addiu - register r[2] now contains 0x00000002
005: subu  - register r[4] now contains 0xfffffffd
006: blez  - branch taken to 0x00000003
003: addu  - register r[1] now contains 0x00000003
004: addiu - register r[2] now contains 0x00000003
005: subu  - register r[4] now contains 0xfffffffe
006: blez  - branch taken to 0x00000003
003: addu  - register r[1] now contains 0x00000006
004: addiu - register r[2] now contains 0x00000004
005: subu  - register r[4] now contains 0xffffffff
006: blez  - branch taken to 0x00000003
003: addu  - register r[1] now contains 0x0000000a
004: addiu - register r[2] now contains 0x00000005
005: subu  - register r[4] now contains 0x00000000
006: blez  - branch taken to 0x00000003
003: addu  - register r[1] now contains 0x0000000f
004: addiu - register r[2] now contains 0x00000006
005: subu  - register r[4] now contains 0x00000001
006: blez  - branch untaken
007: hlt

contents of memory
addr value
000: 24030005
001: 00000821
002: 24020001
003: 00220821
004: 24420001
005: 00432023
006: 1880fffc
007: 00000000

instruction class counts (omits hlt instruction)
  alu ops            18
  loads/stores        0
  jumps/branches      5
total                23

memory access counts (omits hlt instruction)
  inst. fetches      23
  loads               0
  stores              0
total                23

transfer of control counts
  jumps               0
  jump-and-links      0
  taken branches      4
  untaken branches    1
total                 5
```

Guidelines:

You may work individually or in teams of two.

The code should be written totally by yourself or your team of two, but you may discuss the project requirements and the concepts with me or with anyone in the class.

You should not send code to anyone or receive code from anyone, whether by email, printed listings, photos, visual display on a workstation/laptop/cell-phone/etc. screen, or any other method of communication. Do not post the assignment, or a request for help, or your code on any web sites.

The key idea is that you shouldn't short-circuit the learning process for others once you know the answer. (And you shouldn't burden anyone else with inappropriate requests for code or "answers" and thus short-circuit your own learning process.)

Turn in uncompressed source code using handin.cs.clemson.edu before midnight of the due date of Tuesday, Sept. 25. The late penalty is 10% off per day late, up to five days.

Comments are not graded, but you should include your name and your teammate's name if you work in a team. Only one person from a team is required to submit the simulator.

The simulator can be written in any language but must run on school Ubuntu systems and must accept input coming from stdin rather than a named file.
