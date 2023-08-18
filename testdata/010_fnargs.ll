; ModuleID = '010_fnargs.c'
source_filename = "010_fnargs.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

; Function Attrs: noinline nounwind optnone
define dso_local i16 @fn(i16 noundef %0, i16 noundef %1, i16 noundef %2, i16 noundef %3) #0 {
  %5 = alloca i16, align 2
  %6 = alloca i16, align 2
  %7 = alloca i16, align 2
  %8 = alloca i16, align 2
  store i16 %0, i16* %5, align 2
  store i16 %1, i16* %6, align 2
  store i16 %2, i16* %7, align 2
  store i16 %3, i16* %8, align 2
  %9 = load i16, i16* %5, align 2
  %10 = load i16, i16* %6, align 2
  %11 = add nsw i16 %9, %10
  %12 = load i16, i16* %7, align 2
  %13 = add nsw i16 %11, %12
  %14 = load i16, i16* %8, align 2
  %15 = add nsw i16 %13, %14
  ret i16 %15
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  %2 = call i16 @fn(i16 noundef 2, i16 noundef 20, i16 noundef 4, i16 noundef 16)
  ret i16 %2
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
