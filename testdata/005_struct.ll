; ModuleID = '005_struct.c'
source_filename = "005_struct.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

%struct.anon = type { i16, i16 }

@s = dso_local global %struct.anon zeroinitializer, align 2

; Function Attrs: noinline nounwind optnone
define dso_local void @setvars() #0 {
  store i16 3, i16* getelementptr inbounds (%struct.anon, %struct.anon* @s, i32 0, i32 0), align 2
  store i16 5, i16* getelementptr inbounds (%struct.anon, %struct.anon* @s, i32 0, i32 1), align 2
  ret void
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  %1 = alloca i16, align 2
  store i16 0, i16* %1, align 2
  call void @setvars()
  %2 = load i16, i16* getelementptr inbounds (%struct.anon, %struct.anon* @s, i32 0, i32 1), align 2
  %3 = load i16, i16* getelementptr inbounds (%struct.anon, %struct.anon* @s, i32 0, i32 0), align 2
  %4 = sub nsw i16 %2, %3
  %5 = sub nsw i16 %4, 2
  ret i16 %5
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
