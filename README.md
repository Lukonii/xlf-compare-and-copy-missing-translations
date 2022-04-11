# xlf-compare-and-copy-missing-translations
This small script will compare two .xlf (.txt) files, and for every source tag which is not translated in one will try to find translation in other and paste the missing translation.

```xlf
<trans-unit id="xxx" size-unit="char" translate="yes" xml:space="preserve">
  <source>Description</source>
  <target></target>
  <note from="Developer" annotates="general" priority="2"/>
  <note from="Xliff Generator" annotates="general" priority="3"></note>
</trans-unit>
```
<b>`<target></target>`</b> will be replaces with translated one. <b>`<target>Descripción</target>`</b>
</br>It will also work if `<target>` tag is missing by inserting that line.

* donwload main.exe and run it
* Enter the name of 'source.xlf', need-translation file
* Enter the name of 'translation.xlf', translated file
