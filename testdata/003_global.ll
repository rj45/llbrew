; ModuleID = '003_global.c'
source_filename = "003_global.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

@glb = dso_local global i16 39, align 2
@other = dso_local global i16 9, align 2

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  store i16 3, i16* @other, align 2
  %2 = load i16, i16* @glb, align 2
  %3 = load i16, i16* @other, align 2
  %4 = add nsw i16 %2, %3
  ret i16 %4
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
