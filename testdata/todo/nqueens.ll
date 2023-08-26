; ModuleID = '012_nqueens.c'
source_filename = "012_nqueens.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

@t = dso_local global [64 x i16] zeroinitializer, align 2
@N = dso_local global i16 0, align 2

; Function Attrs: noinline nounwind optnone
define dso_local i16 @chk(i16 noundef %0, i16 noundef %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca i16, align 2
  %5 = alloca i16, align 2
  %6 = alloca i16, align 2
  store i16 %0, i16* %3, align 2
  store i16 %1, i16* %4, align 2
  store i16 0, i16* %5, align 2
  store i16 0, i16* %6, align 2
  br label %7

7:                                                ; preds = %131, %2
  %8 = load i16, i16* %5, align 2
  %9 = icmp slt i16 %8, 8
  br i1 %9, label %10, label %134

10:                                               ; preds = %7
  %11 = load i16, i16* %6, align 2
  %12 = load i16, i16* %3, align 2
  %13 = load i16, i16* %5, align 2
  %14 = mul nsw i16 8, %13
  %15 = add nsw i16 %12, %14
  %16 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %15
  %17 = load i16, i16* %16, align 2
  %18 = add nsw i16 %11, %17
  store i16 %18, i16* %6, align 2
  %19 = load i16, i16* %6, align 2
  %20 = load i16, i16* %5, align 2
  %21 = load i16, i16* %4, align 2
  %22 = mul nsw i16 8, %21
  %23 = add nsw i16 %20, %22
  %24 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %23
  %25 = load i16, i16* %24, align 2
  %26 = add nsw i16 %19, %25
  store i16 %26, i16* %6, align 2
  %27 = load i16, i16* %3, align 2
  %28 = load i16, i16* %5, align 2
  %29 = add nsw i16 %27, %28
  %30 = icmp slt i16 %29, 8
  %31 = zext i1 %30 to i16
  %32 = load i16, i16* %4, align 2
  %33 = load i16, i16* %5, align 2
  %34 = add nsw i16 %32, %33
  %35 = icmp slt i16 %34, 8
  %36 = zext i1 %35 to i16
  %37 = and i16 %31, %36
  %38 = icmp ne i16 %37, 0
  br i1 %38, label %39, label %52

39:                                               ; preds = %10
  %40 = load i16, i16* %6, align 2
  %41 = load i16, i16* %3, align 2
  %42 = load i16, i16* %5, align 2
  %43 = add nsw i16 %41, %42
  %44 = load i16, i16* %4, align 2
  %45 = load i16, i16* %5, align 2
  %46 = add nsw i16 %44, %45
  %47 = mul nsw i16 8, %46
  %48 = add nsw i16 %43, %47
  %49 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %48
  %50 = load i16, i16* %49, align 2
  %51 = add nsw i16 %40, %50
  store i16 %51, i16* %6, align 2
  br label %52

52:                                               ; preds = %39, %10
  %53 = load i16, i16* %3, align 2
  %54 = load i16, i16* %5, align 2
  %55 = add nsw i16 %53, %54
  %56 = icmp slt i16 %55, 8
  %57 = zext i1 %56 to i16
  %58 = load i16, i16* %4, align 2
  %59 = load i16, i16* %5, align 2
  %60 = sub nsw i16 %58, %59
  %61 = icmp sge i16 %60, 0
  %62 = zext i1 %61 to i16
  %63 = and i16 %57, %62
  %64 = icmp ne i16 %63, 0
  br i1 %64, label %65, label %78

65:                                               ; preds = %52
  %66 = load i16, i16* %6, align 2
  %67 = load i16, i16* %3, align 2
  %68 = load i16, i16* %5, align 2
  %69 = add nsw i16 %67, %68
  %70 = load i16, i16* %4, align 2
  %71 = load i16, i16* %5, align 2
  %72 = sub nsw i16 %70, %71
  %73 = mul nsw i16 8, %72
  %74 = add nsw i16 %69, %73
  %75 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %74
  %76 = load i16, i16* %75, align 2
  %77 = add nsw i16 %66, %76
  store i16 %77, i16* %6, align 2
  br label %78

78:                                               ; preds = %65, %52
  %79 = load i16, i16* %3, align 2
  %80 = load i16, i16* %5, align 2
  %81 = sub nsw i16 %79, %80
  %82 = icmp sge i16 %81, 0
  %83 = zext i1 %82 to i16
  %84 = load i16, i16* %4, align 2
  %85 = load i16, i16* %5, align 2
  %86 = add nsw i16 %84, %85
  %87 = icmp slt i16 %86, 8
  %88 = zext i1 %87 to i16
  %89 = and i16 %83, %88
  %90 = icmp ne i16 %89, 0
  br i1 %90, label %91, label %104

91:                                               ; preds = %78
  %92 = load i16, i16* %6, align 2
  %93 = load i16, i16* %3, align 2
  %94 = load i16, i16* %5, align 2
  %95 = sub nsw i16 %93, %94
  %96 = load i16, i16* %4, align 2
  %97 = load i16, i16* %5, align 2
  %98 = add nsw i16 %96, %97
  %99 = mul nsw i16 8, %98
  %100 = add nsw i16 %95, %99
  %101 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %100
  %102 = load i16, i16* %101, align 2
  %103 = add nsw i16 %92, %102
  store i16 %103, i16* %6, align 2
  br label %104

104:                                              ; preds = %91, %78
  %105 = load i16, i16* %3, align 2
  %106 = load i16, i16* %5, align 2
  %107 = sub nsw i16 %105, %106
  %108 = icmp sge i16 %107, 0
  %109 = zext i1 %108 to i16
  %110 = load i16, i16* %4, align 2
  %111 = load i16, i16* %5, align 2
  %112 = sub nsw i16 %110, %111
  %113 = icmp sge i16 %112, 0
  %114 = zext i1 %113 to i16
  %115 = and i16 %109, %114
  %116 = icmp ne i16 %115, 0
  br i1 %116, label %117, label %130

117:                                              ; preds = %104
  %118 = load i16, i16* %6, align 2
  %119 = load i16, i16* %3, align 2
  %120 = load i16, i16* %5, align 2
  %121 = sub nsw i16 %119, %120
  %122 = load i16, i16* %4, align 2
  %123 = load i16, i16* %5, align 2
  %124 = sub nsw i16 %122, %123
  %125 = mul nsw i16 8, %124
  %126 = add nsw i16 %121, %125
  %127 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %126
  %128 = load i16, i16* %127, align 2
  %129 = add nsw i16 %118, %128
  store i16 %129, i16* %6, align 2
  br label %130

130:                                              ; preds = %117, %104
  br label %131

131:                                              ; preds = %130
  %132 = load i16, i16* %5, align 2
  %133 = add nsw i16 %132, 1
  store i16 %133, i16* %5, align 2
  br label %7, !llvm.loop !2

134:                                              ; preds = %7
  %135 = load i16, i16* %6, align 2
  ret i16 %135
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @go(i16 noundef %0, i16 noundef %1, i16 noundef %2) #0 {
  %4 = alloca i16, align 2
  %5 = alloca i16, align 2
  %6 = alloca i16, align 2
  %7 = alloca i16, align 2
  store i16 %0, i16* %5, align 2
  store i16 %1, i16* %6, align 2
  store i16 %2, i16* %7, align 2
  %8 = load i16, i16* %5, align 2
  %9 = icmp eq i16 %8, 8
  br i1 %9, label %10, label %13

10:                                               ; preds = %3
  %11 = load i16, i16* @N, align 2
  %12 = add nsw i16 %11, 1
  store i16 %12, i16* @N, align 2
  store i16 0, i16* %4, align 2
  br label %55

13:                                               ; preds = %3
  br label %14

14:                                               ; preds = %51, %13
  %15 = load i16, i16* %7, align 2
  %16 = icmp slt i16 %15, 8
  br i1 %16, label %17, label %54

17:                                               ; preds = %14
  br label %18

18:                                               ; preds = %47, %17
  %19 = load i16, i16* %6, align 2
  %20 = icmp slt i16 %19, 8
  br i1 %20, label %21, label %50

21:                                               ; preds = %18
  %22 = load i16, i16* %6, align 2
  %23 = load i16, i16* %7, align 2
  %24 = call i16 @chk(i16 noundef %22, i16 noundef %23)
  %25 = icmp eq i16 %24, 0
  br i1 %25, label %26, label %46

26:                                               ; preds = %21
  %27 = load i16, i16* %6, align 2
  %28 = load i16, i16* %7, align 2
  %29 = mul nsw i16 8, %28
  %30 = add nsw i16 %27, %29
  %31 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %30
  %32 = load i16, i16* %31, align 2
  %33 = add nsw i16 %32, 1
  store i16 %33, i16* %31, align 2
  %34 = load i16, i16* %5, align 2
  %35 = add nsw i16 %34, 1
  %36 = load i16, i16* %6, align 2
  %37 = load i16, i16* %7, align 2
  %38 = call i16 @go(i16 noundef %35, i16 noundef %36, i16 noundef %37)
  %39 = load i16, i16* %6, align 2
  %40 = load i16, i16* %7, align 2
  %41 = mul nsw i16 8, %40
  %42 = add nsw i16 %39, %41
  %43 = getelementptr inbounds [64 x i16], [64 x i16]* @t, i16 0, i16 %42
  %44 = load i16, i16* %43, align 2
  %45 = add nsw i16 %44, -1
  store i16 %45, i16* %43, align 2
  br label %46

46:                                               ; preds = %26, %21
  br label %47

47:                                               ; preds = %46
  %48 = load i16, i16* %6, align 2
  %49 = add nsw i16 %48, 1
  store i16 %49, i16* %6, align 2
  br label %18, !llvm.loop !4

50:                                               ; preds = %18
  store i16 0, i16* %6, align 2
  br label %51

51:                                               ; preds = %50
  %52 = load i16, i16* %7, align 2
  %53 = add nsw i16 %52, 1
  store i16 %53, i16* %7, align 2
  br label %14, !llvm.loop !5

54:                                               ; preds = %14
  store i16 0, i16* %4, align 2
  br label %55

55:                                               ; preds = %54, %10
  %56 = load i16, i16* %4, align 2
  ret i16 %56
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  %2 = call i16 @go(i16 noundef 0, i16 noundef 0, i16 noundef 0)
  %3 = load i16, i16* @N, align 2
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
