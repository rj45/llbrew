; ModuleID = '008_string.c'
source_filename = "008_string.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

@.str = private unnamed_addr constant [13 x i8] c"hello world!\00", align 1

; Function Attrs: noinline nounwind optnone
define dso_local i16 @len(i8* noundef %0) #0 {
  %2 = alloca i8*, align 2
  %3 = alloca i16, align 2
  store i8* %0, i8** %2, align 2
  store i16 0, i16* %3, align 2
  br label %4

4:                                                ; preds = %14, %1
  %5 = load i8*, i8** %2, align 2
  %6 = load i16, i16* %3, align 2
  %7 = getelementptr inbounds i8, i8* %5, i16 %6
  %8 = load i8, i8* %7, align 1
  %9 = sext i8 %8 to i16
  %10 = icmp eq i16 %9, 0
  br i1 %10, label %11, label %13

11:                                               ; preds = %4
  %12 = load i16, i16* %3, align 2
  ret i16 %12

13:                                               ; preds = %4
  br label %14

14:                                               ; preds = %13
  %15 = load i16, i16* %3, align 2
  %16 = add nsw i16 %15, 1
  store i16 %16, i16* %3, align 2
  br label %4
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  %2 = alloca i8*, align 2
  store i16 0, i16* %1, align 2
  store i8* getelementptr inbounds ([13 x i8], [13 x i8]* @.str, i16 0, i16 0), i8** %2, align 2
  %3 = load i8*, i8** %2, align 2
  %4 = call i16 @len(i8* noundef %3)
  ret i16 %4
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
