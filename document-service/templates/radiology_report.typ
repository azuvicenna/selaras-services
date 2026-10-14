#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[INSTALASI RADIOLOGI & DIAGNOSTIK IMAGING] \
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
      align(left)[#text(7pt, fill: luma(100))[Hasil Ekspertise Radiologi - Terverifikasi Sp.Rad]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 10pt)

#align(center)[
  #text(12pt, weight: "bold")[LEMBAR HASIL PEMERIKSAAN RADIOLOGI] \
  #text(9pt)[Nomor Ekspertise: {{DOC_NUMBER}}]
]

#v(0.5em)

#table(
  columns: (120pt, 8pt, 1fr, 100pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Pemeriksaan], [:], [{{ISSUED_DATE}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Modalitas / Alat], [:], [{{MODALITY}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Dokter Perujuk], [:], [{{DOKTER_PERUJUK}}],
  [Jenis Pemeriksaan], [:], [{{EXAM_NAME}}], [No. Aksesi / DICOM], [:], [{{DICOM_UID}}],
)

#v(0.5em)
#line(length: 100%, stroke: 0.5pt)

*Klinis / Indikasi:* {{INDIKASI_KLINIS}}

#v(0.5em)
*TEMUAN (FINDINGS):*
#block(inset: (left: 10pt))[
  {{FINDINGS}}
]

#v(0.5em)
*KESIMPULAN (IMPRESSION):*
#block(inset: (left: 10pt))[
  *{{IMPRESSION}}*
]

#v(0.5em)
*Saran / Rekomendasi:* {{RECOMMENDATION}}

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI DIGITAL Sp.Rad] \
        #text(6.5pt)[Telah diverifikasi dan ditandatangani secara elektronik:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Spesialis Radiologi (Sp.Rad), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
