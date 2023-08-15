; ModuleID = '006_structptr.c'
source_filename = "006_structptr.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

%struct.S = type { i16, i16 }

; Function Attrs: noinline nounwind optnone
define dso_local void @setS(%struct.S* noundef %0) #0 {
  %2 = alloca %struct.S*, align 2
  store %struct.S* %0, %struct.S** %2, align 2
  %3 = load %struct.S*, %struct.S** %2, align 2
  %4 = getelementptr inbounds %struct.S, %struct.S* %3, i32 0, i32 0
  store i16 1, i16* %4, align 2
  %5 = load %struct.S*, %struct.S** %2, align 2
  %6 = getelementptr inbounds %struct.S, %struct.S* %5, i32 0, i32 1
  store i16 2, i16* %6, align 2
  ret void
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  %2 = alloca %struct.S, align 2
  store i16 0, i16* %1, align 2
  call void @setS(%struct.S* noundef %2)
  %3 = getelementptr inbounds %struct.S, %struct.S* %2, i32 0, i32 1
  %4 = load i16, i16* %3, align 2
  %5 = getelementptr inbounds %struct.S, %struct.S* %2, i32 0, i32 0
  %6 = load i16, i16* %5, align 2
  %7 = add nsw i16 %4, %6
  %8 = sub nsw i16 %7, 3
  ret i16 %8
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
