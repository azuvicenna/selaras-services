#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[INSTALASI REKAM MEDIS & MANAJEMEN INFORMASI KESEHATAN] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Salinan Ringkasan Berkas Rekam Medis - Rahasia Medis]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[RINGKASAN REKAM MEDIS PASIEN (EMR SUMMARY)] \
  #text(9pt)[Nomor Dokumen: {{DOC_NUMBER}}]
]

#v(0.8em)

*PROFIL PASIEN:*
#table(
  columns: (130pt, 8pt, 1fr, 110pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [No. Rekam Medis], [:], [{{PATIENT_ID}}],
  [Nomor NIK / KTP], [:], [{{NIK}}], [Golongan Darah], [:], [{{GOLONGAN_DARAH}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Status Alergi], [:], [*{{ALERGI}}*],
)

#v(0.5em)
*RIWAYAT KUNJUNGAN & EPISODE PELAYANAN:*
#table(
  columns: (80pt, 100pt, 120pt, 1fr),
  stroke: 0.5pt,
  table.header([*Tanggal*], [*Unit Layanan*], [*DPJP / Dokter*], [*Diagnosa Klinis (ICD-10)*]),
  {{RIWAYAT_KUNJUNGAN_ROWS}}
)

#v(0.5em)
*RIWAYAT TINDAKAN & OPERASI:*
#table(
  columns: (80pt, 140pt, 1fr),
  stroke: 0.5pt,
  table.header([*Tanggal*], [*Tindakan / Prosedur (ICD-9)*], [*Dokter Operator*]),
  {{RIWAYAT_TINDAKAN_ROWS}}
)

#v(0.5em)
*DAFTAR PENGOBATAN JANGKA PANJANG (OBAT KRONIS):* \
{{OBAT_KRONIS}}

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI SISTEM ELEKTRONIK] \
        #text(6.5pt)[Dokumen sah tervalidasi dari basis data EMR RS Selaras:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Kepala Instalasi Rekam Medis, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      NIP: {{DOCTOR_SIP}}
    ]
  ]
)
