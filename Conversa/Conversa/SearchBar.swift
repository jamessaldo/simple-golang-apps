//
//  SearchBar.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//

import SwiftUI

struct SearchBar: View {
    @Binding var text: String
    
    @FocusState private var isFocused: Bool
    
    var body: some View {
        HStack {
            TextField("Search ...", text: $text)
                .padding(7)
                .padding(.horizontal, 25)
                .background(Color.white)
                .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                .cornerRadius(10).overlay(
                    HStack {
                        Image(systemName: "magnifyingglass")
                            .foregroundColor(.gray)
                            .frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)
                            .padding(.leading, 8)
                        
                        if isFocused {
                            Button(action: {
                                self.text = ""
                            }) {
                                Image(systemName: "multiply.circle.fill")
                                    .foregroundColor(.gray)
                                    .padding(.trailing, 8)
                            }
                        }
                    }
                )
                .padding(.horizontal, 10)
                .focused($isFocused)
                .onTapGesture {
                    self.isFocused = true
                }
            
            if self.isFocused {
                Button(action: {
                    self.isFocused = false
                    self.text = ""
                    // Dismiss the keyboard
                    UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
                }) {
                    Text("Cancel")
                }
                .padding(.trailing, 10)
                .transition(.move(edge: .trailing))
                .animation(.default)
            }
        }
    }
}

struct SearchBar_Previews: PreviewProvider {
    static var previews: some View {
        SearchBar(text: .constant(""))
    }
}
