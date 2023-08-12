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
  %7 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  %8 = load i16, i16* %3, align 2
  %9 = icmp sle i16 %8, 1
  br i1 %9, label %10, label %12

10:                                               ; preds = %1
  %11 = load i16, i16* %3, align 2
  store i16 %11, i16* %2, align 2
  br label %30

12:                                               ; preds = %1
  store i16 0, i16* %4, align 2
  store i16 1, i16* %5, align 2
  store i16 0, i16* %6, align 2
  store i16 2, i16* %7, align 2
  br label %13

13:                                               ; preds = %23, %12
  %14 = load i16, i16* %7, align 2
  %15 = load i16, i16* %3, align 2
  %16 = icmp slt i16 %14, %15
  br i1 %16, label %17, label %26

17:                                               ; preds = %13
  %18 = load i16, i16* %4, align 2
  %19 = load i16, i16* %5, align 2
  %20 = add nsw i16 %18, %19
  store i16 %20, i16* %6, align 2
  %21 = load i16, i16* %5, align 2
  store i16 %21, i16* %4, align 2
  %22 = load i16, i16* %6, align 2
  store i16 %22, i16* %5, align 2
  br label %23

23:                                               ; preds = %17
  %24 = load i16, i16* %7, align 2
  %25 = add nsw i16 %24, 1
  store i16 %25, i16* %7, align 2
  br label %13, !llvm.loop !2

26:                                               ; preds = %13
  %27 = load i16, i16* %5, align 2
  %28 = load i16, i16* %4, align 2
  %29 = add nsw i16 %27, %28
  store i16 %29, i16* %2, align 2
  br label %30

30:                                               ; preds = %26, %10
  %31 = load i16, i16* %2, align 2
  ret i16 %31
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
