; ModuleID = '007_recurptr.c'
source_filename = "007_recurptr.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

%struct.S = type { %struct.S*, i16 }

@s = dso_local global %struct.S zeroinitializer, align 2

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  store i16 42, i16* getelementptr inbounds (%struct.S, %struct.S* @s, i32 0, i32 1), align 2
  store %struct.S* @s, %struct.S** getelementptr inbounds (%struct.S, %struct.S* @s, i32 0, i32 0), align 2
  %2 = load %struct.S*, %struct.S** getelementptr inbounds (%struct.S, %struct.S* @s, i32 0, i32 0), align 2
  %3 = getelementptr inbounds %struct.S, %struct.S* %2, i32 0, i32 0
  %4 = load %struct.S*, %struct.S** %3, align 2
  %5 = getelementptr inbounds %struct.S, %struct.S* %4, i32 0, i32 0
  %6 = load %struct.S*, %struct.S** %5, align 2
  %7 = getelementptr inbounds %struct.S, %struct.S* %6, i32 0, i32 0
  %8 = load %struct.S*, %struct.S** %7, align 2
  %9 = getelementptr inbounds %struct.S, %struct.S* %8, i32 0, i32 0
  %10 = load %struct.S*, %struct.S** %9, align 2
  %11 = getelementptr inbounds %struct.S, %struct.S* %10, i32 0, i32 1
  %12 = load i16, i16* %11, align 2
  ret i16 %12
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
