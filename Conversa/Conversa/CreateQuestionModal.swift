//
//  CreateQuestionModal.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//

import SwiftUI

struct CreateQuestionModal: View {
    @Binding var isPresented: Bool
    @Binding var title: String
    @Binding var description: String
    
    var body: some View {
        NavigationView{
            VStack(spacing: 16) {
                Text("Title")
                    .font(.subheadline)
                    .bold()
                    .foregroundColor(Color(hex: "#383838"))
                    .padding(.top, 10)
                TextEditor(text: $title)
                    .frame(width: 284, height: 38)
                    .padding(.horizontal, 16)
                    .foregroundColor(Color(hex: "#383838"))
                    .background(Color(hex: "#F4F4F4"))
                    .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                    .cornerRadius(10)
                
                Text("Description")
                    .font(.subheadline)
                    .bold()
                    .foregroundColor(Color(hex: "#383838"))
                    .padding(.top, 10)
                TextEditor(text: $description)
                    .frame(width: 284, height: 159)
                    .padding(.horizontal, 16)
                    .foregroundColor(Color(hex: "#383838"))
                    .background(Color(hex: "#F4F4F4"))
                    .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                    .cornerRadius(10)
                
                Button(action: {
                    isPresented.toggle()
                    title = ""
                    description = ""
                }, label: {
                    Text("Post")
                        .foregroundColor(Color(hex: "#383838"))
                        .frame(
                            minWidth: 0,
                            maxWidth: 314,
                            minHeight: 0,
                            maxHeight: 38
                        )
                        .background(Color(hex: "#FFEBBB"))
                        .cornerRadius(10)
                        .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                })
            }
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: {
                        isPresented.toggle()
                        title = ""
                        description = ""}
                    ) {
                        Text("Cancel")
                            .foregroundColor(Color(hex: "#383838"))
                            .frame(
                                width: 64 , height: 30, alignment: .center
                            )
                            .background(Color(hex: "#FFEBBB"))
                            .cornerRadius(10)
                            .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                    }
                }
            }
        }
    }
}
