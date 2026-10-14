#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[INSTALASI DIAGNOSTIK KARDIO-NEUROLOGI] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Laporan Hasil Elektromedik (EKG / EEG) - Terverifikasi Sp.JP / Sp.N]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[LAPORAN HASIL REKAM ELEKTROMEDIK (EKG / EEG)] \
  #text(9pt)[Nomor Laporan: {{DOC_NUMBER}}]
]

#v(0.8em)

#table(
  columns: (120pt, 8pt, 1fr, 110pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Pemeriksaan], [:], [{{ISSUED_DATE}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Jenis Pemeriksaan], [:], [{{JENIS_TES}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Dokter Perujuk], [:], [{{DOKTER_PERUJUK}}],
  [Teknisi Elektromedik], [:], [{{TEKNISI_NAME}}], [Indikasi Klinis], [:], [{{INDIKASI}}],
)

#v(0.5em)
#line(length: 100%, stroke: 0.5pt)

*PARAMETER PENGUKURAN ELEKTROMEDIK:*
#table(
  columns: (1fr, 1fr, 1fr),
  stroke: 0.5pt,
  [ *Irama (Rhythm):* \ {{IRAMA}} ],
  [ *Laju Detak / Frekuensi:* \ {{HEART_RATE}} bpm / Hz ],
  [ *Axis:* \ {{AXIS}} ],
  [ *Interval PR:* \ {{PR_INTERVAL}} ms ],
  [ *Durasi Kompleks QRS:* \ {{QRS_DURATION}} ms ],
  [ *Interval QT / QTc:* \ {{QT_INTERVAL}} ms ],
)

#v(0.5em)
*DESKRIPSI GELOMBANG & TEMUAN ELEKTROFISIOLOGI:*
#block(inset: (left: 10pt))[
  {{TEMUAN_GELOMBANG}}
]

#v(0.5em)
*KESIMPULAN & INTERPRETASI KLINIS:*
#block(inset: (left: 10pt))[
  *{{KESIMPULAN}}*
]

#v(0.5em)
*Saran / Tindak Lanjut:* {{SARAN}}

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI DIGITAL DOKTER SPESIALIS] \
        #text(6.5pt)[Telah diverifikasi dan ditandatangani secara elektronik:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Pemeriksa (Sp.JP / Sp.N), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
