; ModuleID = '011_loopinter.c'
source_filename = "011_loopinter.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

@A = dso_local global [64 x [64 x i16]] zeroinitializer, align 2
@y = dso_local global i16 0, align 2

; Function Attrs: noinline nounwind optnone
define dso_local void @init() #0 {
  %1 = alloca i16, align 2
  %2 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  br label %3

3:                                                ; preds = %22, %0
  %4 = load i16, i16* %1, align 2
  %5 = icmp ult i16 %4, 64
  br i1 %5, label %6, label %25

6:                                                ; preds = %3
  store i16 0, i16* %2, align 2
  br label %7

7:                                                ; preds = %18, %6
  %8 = load i16, i16* %2, align 2
  %9 = icmp ult i16 %8, 64
  br i1 %9, label %10, label %21

10:                                               ; preds = %7
  %11 = load i16, i16* %1, align 2
  %12 = load i16, i16* %2, align 2
  %13 = add i16 %11, %12
  %14 = load i16, i16* %1, align 2
  %15 = getelementptr inbounds [64 x [64 x i16]], [64 x [64 x i16]]* @A, i16 0, i16 %14
  %16 = load i16, i16* %2, align 2
  %17 = getelementptr inbounds [64 x i16], [64 x i16]* %15, i16 0, i16 %16
  store i16 %13, i16* %17, align 2
  br label %18

18:                                               ; preds = %10
  %19 = load i16, i16* %2, align 2
  %20 = add i16 %19, 1
  store i16 %20, i16* %2, align 2
  br label %7, !llvm.loop !2

21:                                               ; preds = %7
  br label %22

22:                                               ; preds = %21
  %23 = load i16, i16* %1, align 2
  %24 = add i16 %23, 1
  store i16 %24, i16* %1, align 2
  br label %3, !llvm.loop !4

25:                                               ; preds = %3
  ret void
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @interchange() #0 {
  %1 = alloca i16, align 2
  %2 = alloca i16, align 2
  %3 = alloca i16, align 2
  store i16 0, i16* %2, align 2
  br label %4

4:                                                ; preds = %29, %0
  %5 = load i16, i16* %2, align 2
  %6 = icmp ult i16 %5, 64
  br i1 %6, label %7, label %32

7:                                                ; preds = %4
  store i16 0, i16* %1, align 2
  store i16 0, i16* %3, align 2
  br label %8

8:                                                ; preds = %25, %7
  %9 = load i16, i16* %3, align 2
  %10 = icmp ult i16 %9, 64
  br i1 %10, label %11, label %28

11:                                               ; preds = %8
  %12 = load i16, i16* %2, align 2
  %13 = getelementptr inbounds [64 x [64 x i16]], [64 x [64 x i16]]* @A, i16 0, i16 %12
  %14 = load i16, i16* %3, align 2
  %15 = getelementptr inbounds [64 x i16], [64 x i16]* %13, i16 0, i16 %14
  %16 = load i16, i16* %15, align 2
  %17 = add i16 %16, 1
  store i16 %17, i16* %15, align 2
  %18 = load i16, i16* %2, align 2
  %19 = getelementptr inbounds [64 x [64 x i16]], [64 x [64 x i16]]* @A, i16 0, i16 %18
  %20 = load i16, i16* %3, align 2
  %21 = getelementptr inbounds [64 x i16], [64 x i16]* %19, i16 0, i16 %20
  %22 = load i16, i16* %21, align 2
  %23 = load i16, i16* %1, align 2
  %24 = add i16 %23, %22
  store i16 %24, i16* %1, align 2
  br label %25

25:                                               ; preds = %11
  %26 = load i16, i16* %3, align 2
  %27 = add i16 %26, 1
  store i16 %27, i16* %3, align 2
  br label %8, !llvm.loop !5

28:                                               ; preds = %8
  br label %29

29:                                               ; preds = %28
  %30 = load i16, i16* %2, align 2
  %31 = add i16 %30, 1
  store i16 %31, i16* %2, align 2
  br label %4, !llvm.loop !6

32:                                               ; preds = %4
  %33 = load i16, i16* %1, align 2
  ret i16 %33
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  call void @init()
  %2 = call i16 @interchange()
  %3 = and i16 %2, 127
  ret i16 %3
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
!2 = distinct !{!2, !3}
!3 = !{!"llvm.loop.mustprogress"}
!4 = distinct !{!4, !3}
!5 = distinct !{!5, !3}
!6 = distinct !{!6, !3}
