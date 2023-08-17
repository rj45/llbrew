; ModuleID = '009_summary.c'
source_filename = "009_summary.c"
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
define dso_local i16 @sub(i16 noundef %0, i16 noundef %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  store i16 %1, i16* %4, align 2
  %5 = load i16, i16* %3, align 2
  %6 = load i16, i16* %4, align 2
  %7 = sub nsw i16 %5, %6
  ret i16 %7
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @and(i16 noundef %0, i16 noundef %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  store i16 %1, i16* %4, align 2
  %5 = load i16, i16* %3, align 2
  %6 = load i16, i16* %4, align 2
  %7 = and i16 %5, %6
  ret i16 %7
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @or(i16 noundef %0, i16 noundef %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  store i16 %1, i16* %4, align 2
  %5 = load i16, i16* %3, align 2
  %6 = load i16, i16* %4, align 2
  %7 = or i16 %5, %6
  ret i16 %7
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @xor(i16 noundef %0, i16 noundef %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  store i16 %1, i16* %4, align 2
  %5 = load i16, i16* %3, align 2
  %6 = load i16, i16* %4, align 2
  %7 = xor i16 %5, %6
  ret i16 %7
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  %2 = alloca i16, align 2
  %3 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  %4 = call i16 @add(i16 noundef 123, i16 noundef 456)
  store i16 %4, i16* %2, align 2
  %5 = call i16 @sub(i16 noundef 654, i16 noundef 321)
  store i16 %5, i16* %3, align 2
  %6 = load i16, i16* %2, align 2
  %7 = load i16, i16* %3, align 2
  %8 = call i16 @xor(i16 noundef %6, i16 noundef %7)
  store i16 %8, i16* %2, align 2
  %9 = load i16, i16* %2, align 2
  %10 = load i16, i16* %3, align 2
  %11 = call i16 @or(i16 noundef %9, i16 noundef %10)
  store i16 %11, i16* %3, align 2
  %12 = load i16, i16* %3, align 2
  %13 = call i16 @and(i16 noundef %12, i16 noundef 481)
  store i16 %13, i16* %2, align 2
  %14 = load i16, i16* %2, align 2
  %15 = sub nsw i16 %14, 279
  ret i16 %15
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
