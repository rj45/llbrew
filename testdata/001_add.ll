; ModuleID = '001_add.c'
source_filename = "001_add.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

; Function Attrs: noinline nounwind optnone
define dso_local i16 @add(i16 noundef %0, i16 noundef %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  store i16 %1, i16* %4, align 2
  %5 = load i16, i16* %3, align 2
  %6 = load i16, i16* %4, align 2
  %7 = add nsw i16 %5, %6
  ret i16 %7
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  %2 = call i16 @add(i16 noundef 38, i16 noundef 4) #1
  store i16 %2, i16* %1, align 2
  %3 = load i16, i16* %1, align 2
  ret i16 %3
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-builtins" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }
attributes #1 = { nobuiltin "no-builtins" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
