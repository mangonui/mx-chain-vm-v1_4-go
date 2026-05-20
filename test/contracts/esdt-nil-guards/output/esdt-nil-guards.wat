(module
  (type (;0;) (func (param i32 i32 i32 i64 i32) (result i32)))
  (type (;1;) (func (param i32 i32 i32 i64 i32 i32 i32 i32 i32 i32 i32 i32) (result i32)))
  (type (;2;) (func (param i32 i32) (result i32)))
  (type (;3;) (func (param i32 i32 i64 i32)))
  (type (;4;) (func (param i32 i32 i64 i32 i32 i32 i32 i32 i32 i32 i32)))
  (type (;5;) (func (param i64) (result i32)))
  (type (;6;) (func (param i32 i32 i32 i64 i32)))
  (type (;7;) (func))
  (import "env" "getESDTBalance" (func (;0;) (type 0)))
  (import "env" "getESDTTokenData" (func (;1;) (type 1)))
  (import "env" "mBufferNewFromBytes" (func (;2;) (type 2)))
  (import "env" "managedGetESDTBalance" (func (;3;) (type 3)))
  (import "env" "managedGetESDTTokenData" (func (;4;) (type 4)))
  (import "env" "bigIntNew" (func (;5;) (type 5)))
  (import "env" "bigIntGetESDTExternalBalance" (func (;6;) (type 6)))
  (func (;7;) (type 7))
  (func (;8;) (type 7)
    i32.const 0
    i32.const 32
    i32.const 5
    i64.const 0
    i32.const 64
    call 0
    drop)
  (func (;9;) (type 7)
    i32.const 0
    i32.const 32
    i32.const 5
    i64.const 0
    i32.const 1
    i32.const 64
    i32.const 96
    i32.const 128
    i32.const 160
    i32.const 192
    i32.const 2
    i32.const 224
    call 1
    drop)
  (func (;10;) (type 7)
    (local i32 i32)
    i32.const 0
    i32.const 32
    call 2
    local.set 0
    i32.const 32
    i32.const 5
    call 2
    local.set 1
    local.get 0
    local.get 1
    i64.const 0
    i32.const 3
    call 3)
  (func (;11;) (type 7)
    (local i32 i32)
    i32.const 0
    i32.const 32
    call 2
    local.set 0
    i32.const 32
    i32.const 5
    call 2
    local.set 1
    local.get 0
    local.get 1
    i64.const 0
    i32.const 4
    i32.const 5
    i32.const 6
    i32.const 7
    i32.const 8
    i32.const 9
    i32.const 10
    i32.const 11
    call 4)
  (func (;12;) (type 7)
    i64.const 123
    call 5
    drop
    i32.const 0
    i32.const 32
    i32.const 5
    i64.const 0
    i32.const 12
    call 6)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66576))
  (export "memory" (memory 0))
  (export "init" (func 7))
  (export "legacyBalance" (func 8))
  (export "legacyTokenData" (func 9))
  (export "managedBalance" (func 10))
  (export "managedTokenData" (func 11))
  (export "bigIntExternalBalance" (func 12))
  (data (;0;) (i32.const 0) "\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01\01TOKEN"))
