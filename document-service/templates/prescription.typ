#set page(
  paper: "a5",
  margin: (x: 1.5cm, y: 1.8cm),
  header: [
    #align(center)[
      #text(11pt, weight: "bold")[INSTALASI FARMASI RS SELARAS MEDIKA] \
      #text(7pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(1pt)
      #line(length: 100%, stroke: 1pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #align(center)[#text(6.5pt, fill: luma(100))[Salinan Resep Sah Farmasi RS Selaras Medika]]
  ]
)

#set text(size: 8.5pt)

#align(center)[
  #text(11pt, weight: "bold")[SALINAN RESEP DOKTER] \
  #text(8pt)[Nomor Resep: {{DOC_NUMBER}}]
]

#v(0.3em)

#table(
  columns: (80pt, 5pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}} ({{PATIENT_GENDER}}, {{PATIENT_DOB}})],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}],
  [Dokter Penulis], [:], [{{DOCTOR_NAME}}],
)

#v(0.5em)
#line(length: 100%, stroke: 0.5pt)
#v(0.3em)

#text(12pt, weight: "bold")[R /]

#block(inset: (left: 15pt))[
  {{PRESCRIPTION_ITEMS}}
]

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 4pt)[
      #align(center)[
        #text(6.5pt, weight: "bold")[VERIFIKASI RESEP] \
        #text(6pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Penulis Resep, \
      #v(2.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
