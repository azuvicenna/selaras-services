#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[INSTALASI LABORATORIUM KLINIK] \
      #text(10pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Hasil Laboratorium Terverifikasi Sp.PK]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[LEMBAR HASIL UJI LABORATORIUM] \
  #text(9pt)[Nomor Laboratorium: {{DOC_NUMBER}}]
]

#v(0.5em)

#table(
  columns: (120pt, 8pt, 1fr, 100pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Pengambilan], [:], [{{TGL_PENGAMBILAN}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Tgl Selesai Uji], [:], [{{ISSUED_DATE}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Dokter Pengirim], [:], [{{DOKTER_PENGIRIM}}],
  [Spesimen / Sampel], [:], [{{SPESIMEN}}], [Analis Lab], [:], [{{ANALIS_NAME}}],
)

#v(0.5em)

#table(
  columns: (1.5fr, 1fr, 1fr, 1fr, 0.8fr),
  stroke: 0.5pt,
  table.header(
    [*Pemeriksaan*], [*Hasil*], [*Nilai Rujukan*], [*Satuan*], [*Status*]
  ),
  {{LAB_ROWS}}
)

#v(0.5em)
*Catatan / Interpretasi:* {{LAB_NOTES}}

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI DIGITAL Sp.PK] \
        #text(6.5pt)[Telah diverifikasi dan ditandatangani secara elektronik:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Penanggung Jawab Laboratorium (Sp.PK), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
