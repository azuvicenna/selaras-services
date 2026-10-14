#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[INSTALASI PATOLOGI ANATOMI & HISTOPATOLOGI] \
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
      align(left)[#text(7pt, fill: luma(100))[Laporan Hasil Patologi Anatomi - Terverifikasi Sp.PA]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[LAPORAN PEMERIKSAAN PATOLOGI ANATOMI (BIOPSI)] \
  #text(9pt)[Nomor Laboratorium PA: {{DOC_NUMBER}}]
]

#v(0.8em)

#table(
  columns: (120pt, 8pt, 1fr, 110pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Terima Jaringan], [:], [{{TGL_TERIMA}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Tgl Selesai Baca], [:], [{{ISSUED_DATE}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Dokter Operator/Pengirim], [:], [{{DOKTER_PENGIRIM}}],
  [Asal Jaringan / Organ], [:], [{{ASAL_JARINGAN}}], [Metode Pengambilan], [:], [{{METODE_BIOPSI}}],
)

#v(0.5em)
#line(length: 100%, stroke: 0.5pt)

*I. PEMERIKSAAN MAKROSKOPIS:*
#block(inset: (left: 10pt))[
  {{MAKROSKOPIS}}
]

#v(0.5em)
*II. PEMERIKSAAN MIKROSKOPIS:*
#block(inset: (left: 10pt))[
  {{MIKROSKOPIS}}
]

#v(0.5em)
*III. KESIMPULAN PATOLOGIS / DIAGNOSA HISTOPATOLOGI:*
#block(inset: (left: 10pt))[
  #text(weight: "bold")[{{KESIMPULAN_PA}}] \
  #text(8.5pt)[Topografi / Morfologi (ICD-O): {{ICD_O_CODE}}]
]

#v(0.5em)
*Catatan / Imunohistokimia (IHC):* {{CATATAN_IHC}}

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI DIGITAL Sp.PA] \
        #text(6.5pt)[Telah diverifikasi dan ditandatangani secara elektronik:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Spesialis Patologi Anatomi (Sp.PA), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
