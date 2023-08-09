; ModuleID = '002_fib.c'
source_filename = "002_fib.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

; Function Attrs: noinline nounwind optnone
define dso_local i16 @fibonacci(i16 noundef %0) #0 {
  %2 = alloca i16, align 2
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  %5 = alloca i16, align 2
  %6 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  %7 = load i16, i16* %3, align 2
  %8 = icmp sle i16 %7, 1
  br i1 %8, label %9, label %11

9:                                                ; preds = %1
  %10 = load i16, i16* %3, align 2
  store i16 %10, i16* %2, align 2
  br label %28

11:                                               ; preds = %1
  store i16 0, i16* %4, align 2
  store i16 1, i16* %5, align 2
  store i16 2, i16* %6, align 2
  br label %12

12:                                               ; preds = %21, %11
  %13 = load i16, i16* %6, align 2
  %14 = load i16, i16* %3, align 2
  %15 = icmp slt i16 %13, %14
  br i1 %15, label %16, label %24

16:                                               ; preds = %12
  %17 = load i16, i16* %5, align 2
  store i16 %17, i16* %4, align 2
  %18 = load i16, i16* %5, align 2
  %19 = load i16, i16* %4, align 2
  %20 = add nsw i16 %18, %19
  store i16 %20, i16* %5, align 2
  br label %21

21:                                               ; preds = %16
  %22 = load i16, i16* %6, align 2
  %23 = add nsw i16 %22, 1
  store i16 %23, i16* %6, align 2
  br label %12, !llvm.loop !2

24:                                               ; preds = %12
  %25 = load i16, i16* %4, align 2
  %26 = load i16, i16* %5, align 2
  %27 = add nsw i16 %25, %26
  store i16 %27, i16* %2, align 2
  br label %28

28:                                               ; preds = %24, %9
  %29 = load i16, i16* %2, align 2
  ret i16 %29
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = call i16 @fibonacci(i16 noundef 7) #1
  ret i16 %1
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-builtins" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }
attributes #1 = { nobuiltin "no-builtins" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
!2 = distinct !{!2, !3}
!3 = !{!"llvm.loop.mustprogress"}
