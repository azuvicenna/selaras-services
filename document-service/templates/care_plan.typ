#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[TIM ASUHAN INTERDISIPLIN / MULTIDISIPLIN] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Rencana Asuhan Terintegrasi (Integrated Clinical Care Plan)]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9pt)

#align(center)[
  #text(12pt, weight: "bold")[RENCANA PERAWATAN KLINIS TERINTEGRASI (CARE PLAN)] \
  #text(8.5pt)[Nomor Dokumen: {{DOC_NUMBER}}]
]

#v(0.5em)

#table(
  columns: (120pt, 8pt, 1fr, 100pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Masuk Rawat], [:], [{{TGL_MASUK}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Target Hari Rawat], [:], [{{TARGET_HARI}} Hari],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Dokter DPJP Utama], [:], [{{DOCTOR_NAME}}],
  [Diagnosa Kerja], [:], [{{DIAGNOSIS}}], [Ruang / Bed], [:], [{{RUANG_RAWAT}}],
)

#v(0.5em)
*TARGET KLINIS / EXPECTED OUTCOME:*
#block(inset: (left: 10pt))[
  {{TARGET_KLINIS}}
]

#v(0.5em)
*RENCANA INTERVENSI MULTIDISIPLIN:*
#table(
  columns: (130pt, 1fr, 120pt),
  stroke: 0.5pt,
  table.header([*Disiplin Profesional (PPA)*], [*Rencana Intervensi & Sasaran Terapi*], [*Penanggung Jawab*]),
  [Dokter DPJP / Spesialis], [{{RENCANA_MEDIS}}], [{{DOCTOR_NAME}}],
  [Keperawatan], [{{RENCANA_KEPERAWATAN}}], [{{NURSE_NAME}}],
  [Farmasi Klinis], [{{RENCANA_FARMASI}}], [{{APOTEKER_NAME}}],
  [Gizi & Nutrisi], [{{RENCANA_GIZI}}], [{{NUTRISIONIS_NAME}}],
  [Rehabilitasi / Fisioterapi], [{{RENCANA_FISIOTERAPI}}], [{{FISIOTERAPIS_NAME}}],
)

#v(0.5em)
*RENCANA EDUKASI PASIEN & KELUARGA:*
#block(inset: (left: 10pt))[
  {{RENCANA_EDUKASI}}
]

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI INTEGRATED CARE] \
        #text(6.5pt)[Rencana terintegrasi tervalidasi elektronik:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Penanggung Jawab Pelayanan (DPJP), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
