; ModuleID = '013_putc.c'
source_filename = "013_putc.c"
target datalayout = "e-m:e-p:16:16-i32:16-i64:16-f32:16-f64:16-a:8-n8:16-S16"
target triple = "msp430-unknown-unknown-elf"

; Function Attrs: noinline nounwind optnone
define dso_local void @putc(i8 noundef signext %0) #0 {
  %2 = alloca i8, align 1
  store i8 %0, i8* %2, align 1
  %3 = load i8, i8* %2, align 1
  %4 = sext i8 %3 to i16
  store i16 %4, i16* inttoptr (i16 -256 to i16*), align 2
  ret void
}

; Function Attrs: noinline nounwind optnone
define dso_local i16 @main() #0 {
  call void @putc(i8 noundef signext 72)
  call void @putc(i8 noundef signext 101)
  call void @putc(i8 noundef signext 108)
  call void @putc(i8 noundef signext 108)
  call void @putc(i8 noundef signext 111)
  call void @putc(i8 noundef signext 114)
  call void @putc(i8 noundef signext 108)
  call void @putc(i8 noundef signext 100)
  call void @putc(i8 noundef signext 33)
  call void @putc(i8 noundef signext 13)
  call void @putc(i8 noundef signext 10)
  ret i16 0
}

attributes #0 = { noinline nounwind optnone "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 2}
!1 = !{!"Apple clang version 14.0.3 (clang-1403.0.22.14.1)"}
