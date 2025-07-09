#!/usr/bin/env perl
use strict;
use warnings;

# Ensure filename is given
die "Usage: $0 <markdown-file>\n" unless @ARGV == 1;
my $filename = $ARGV[0];

# Read the entire file into @lines
open(my $fh, '<', $filename) or die "Could not open '$filename' for reading: $!\n";
my @lines = <$fh>;
close($fh);

my @output;
my $i = 0;

while ($i < @lines) {
    my $line = $lines[$i];

    if ($line =~ /^```([\w+-]+):(.+)\s*$/) {
        my ($lang, $filepath) = ($1, $2);

        # Ensure file exists
        die "Error: File '$filepath' not found\n" unless -e $filepath;

        # Open the included file and read contents
        open(my $code_fh, '<', $filepath) or die "Cannot read '$filepath': $!\n";
        my @code_lines = <$code_fh>;
        close($code_fh);

        # Insert code block
        push @output, "```$lang:$filepath\n";
        push @output, @code_lines;
        chomp($output[-1]);  # Prevent extra newline before closing ```
        push @output, "\n```\n";

        # Skip to closing ```
        $i++;
        while ($i < @lines && $lines[$i] !~ /^```/) {
            $i++;
        }
        $i++;  # Skip the closing ```
    } else {
        push @output, $line;
        $i++;
    }
}

# Write the modified content back to the original file
open(my $out_fh, '>', $filename) or die "Could not open '$filename' for writing: $!\n";
print $out_fh @output;
close($out_fh);
